// SWARM Distributed Monitoring System
// https://github.com/jbuchbinder/swarm-monitor
//
// vim: tabstop=2:softtabstop=2:shiftwidth=2:noexpandtab

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jbuchbinder/swarm-monitor/config"
	"github.com/jbuchbinder/swarm-monitor/logging"
	"github.com/jbuchbinder/swarm-monitor/util"
)

var (
	// ControlThreadRunning says whether the control thread is currently active
	ControlThreadRunning = false
)

// Attempt to use SETNX to assert that we have the control thread.
// If this happens, we have to set the expiration time.
func grabControlThread(c *redis.Client) bool {
	val, err := c.SetNX(ctx, CONTROL_THREAD_LOCK, []byte(fmt.Sprint(config.Config.HostID)), CONTROL_THREAD_EXPIRY*time.Second).Result()
	if err == nil {
		if val {
			// New lock acquired, attempt to set expiry properly
			logging.Infof("grabControlThread: Acquired lock")
			res := c.Expire(ctx, CONTROL_THREAD_LOCK, CONTROL_THREAD_EXPIRY)
			if res.Err() != nil {
				log.Printf("ERR: %s", res.Err().Error())
				return false
			}
			return true
		}
		// Already there
		return false
	}
	log.Printf("ERR: grabControlThread: " + err.Error())
	return true
}

func dropControlThread(c *redis.Client) {
	val, err := c.Get(ctx, CONTROL_THREAD_LOCK).Result()
	if err != nil {
		log.Printf("ERR: " + err.Error())
		return
	}
	if bytes.Compare([]byte(val), []byte(fmt.Sprint(config.Config.HostID))) == 0 {
		logging.Infof("Found control thread with matching id %d, dropping", config.Config.HostID)
		res := c.Del(ctx, CONTROL_THREAD_LOCK)
		if res.Err() != nil {
			log.Printf("ERR: " + res.Err().Error())
			return
		}
	}
}

func extendControlExpiry(c *redis.Client) {
	c.Expire(ctx, CONTROL_THREAD_LOCK, CONTROL_THREAD_EXPIRY*time.Second)
}

func threadControl() {
	if util.ShuttingDown {
		return
	}

	if ControlThreadRunning {
		logging.Debugf("Control thread start attempting, but it looks like it's already running")
		return
	}

	logging.Infof("Starting control thread")

	util.RunningProcesses.Add(1)
	defer util.RunningProcesses.Done()

	c := redis.NewClient(getConnection(REDIS_READWRITE).connspec)

	// Check to see if we need to run control thread, or if it is currently
	// running on another host
	for {
		logging.Debugf("Attempting to grab control thread")
		if grabControlThread(c) {
			// If we're actually starting up, set global running flag
			ControlThreadRunning = true
			for {
				if !ControlThreadRunning {
					// Catch shutdown
					logging.Warnf("ControlThreadRunning was set to false, shutting down control thread")
					return
				}

				if util.ShuttingDown {
					logging.Infof("ControlThread relinquishing control and terminating")
					dropControlThread(c)
					return
				}

				// Endlessly attempt to schedule checks
				logging.Debugf("Attempt to schedule checks")
				members, err := c.SMembers(ctx, CHECKS_LIST).Result()
				if err != nil {
					log.Printf("ERR: Unable to pull from key " + CHECKS_LIST)
				} else {
					logging.Debugf("Members: %#v", members)
					// Pull list of all hosts/services
					for i := 0; i < len(members); i++ {
						logging.Debugf("Scheduling for member %s", members[i])
						// Pull last run and schedule interval to see if this needs to
						// be scheduled for another run, and push onto POLL_QUEUE.
						intervalRaw, err := c.HGet(ctx, string(members[i]), "interval").Result()
						interval, _ := strconv.ParseUint(string(intervalRaw), 10, 64)
						if interval == 0 {
							interval = config.Config.DefaultCheckInterval
						}

						command, err := c.HGet(ctx, string(members[i]), "command").Result()

						typeRaw, err := c.HGet(ctx, string(members[i]), "type").Result()
						checkType, _ := strconv.ParseUint(string(typeRaw), 10, 32)

						curtime := uint64(time.Now().Unix())

						items, err := c.HGetAll(ctx, string(members[i])+":hosts").Result()
						if err == nil {
							for host, v := range items {
								lastrun, _ := strconv.ParseUint(string(v), 10, 64)
								logging.Debugf("curtime = %d, lastrun = %d, diff = %d, interval = %d", curtime, lastrun, curtime-lastrun, interval)
								if curtime-lastrun < interval {
									continue
								}
								go func(c *redis.Client, host string, m string, checkType uint64, command string, curtime uint64) {
									sleep := rand.Intn(int(interval / 2))
									logging.Debugf("Check sleeping for %d seconds", sleep)
									util.RunningProcesses.Add(1)
									defer util.RunningProcesses.Done()
									for i := 0; i <= sleep; i++ {
										time.Sleep(time.Second)
										if util.ShuttingDown {
											return
										}
									}
									logging.Infof("Adding %s : %s to poll queue", m, host)
									// Set lastrun to current time
									_, e := c.HSet(ctx, string(m)+":hosts", host, []byte(fmt.Sprint(curtime))).Result()
									if e != nil {
										log.Printf("ERR: " + e.Error())
									}
									// Also update reverse index
									_, e = c.HSet(ctx, host+":checks", string(m), []byte(fmt.Sprint(curtime))).Result()
									if e != nil {
										log.Printf("ERR: " + e.Error())
									}
									// Form JSON object to serialize onto the scheduler stack
									obj := PollCheck{
										Host:        host,
										CheckName:   string(m),
										EnqueueTime: curtime,
										Type:        uint(checkType),
										Command:     string(command),
									}
									logging.Debugf("Adding to scheduler stack: %#v", obj)
									o, err := json.Marshal(obj)
									if err == nil {
										_, e = c.RPush(ctx, POLL_QUEUE, o).Result()
										if e != nil {
											log.Printf("ERR: " + e.Error())
										}
									} else {
										log.Printf("ERR: " + e.Error())
									}
								}(c, host, members[i], checkType, command, curtime)
							}
						}
					}

				}

				if util.ShuttingDown {
					logging.Infof("ControlThread relinquishing control and terminating")
					dropControlThread(c)
					return
				}

				// Extend control thread expiry
				extendControlExpiry(c)

				// Sleep for a few seconds to avoid CPU piling.
				time.Sleep(2000 * time.Millisecond)
			}
		} else {
			logging.Debugf("ControlThread: already owned, waiting to start.")
		}

		// Sleep for a few seconds to avoid CPU piling.
		for i := 0; i < 15; i++ {
			if util.ShuttingDown {
				logging.Infof("ControlThread relinquishing control and terminating")
				dropControlThread(c)
				return
			}
			time.Sleep(1 * time.Second)
		}
	} // for
}
