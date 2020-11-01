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
	"strconv"
	"time"

	redis "github.com/jbuchbinder/go-redis"
	"github.com/jbuchbinder/swarm-monitor/config"
	"github.com/jbuchbinder/swarm-monitor/util"
)

var (
	ControlThreadRunning = false
)

// Attempt to use SETNX to assert that we have the control thread.
// If this happens, we have to set the expiration time.
func grabControlThread(c redis.Client) bool {
	val, err := c.Setnx(CONTROL_THREAD_LOCK, []byte(fmt.Sprint(config.Config.HostID)))
	if err == nil {
		if val {
			// New lock acquired, attempt to set expiry properly
			log.Printf("INFO: grabControlThread: Acquired lock")
			_, err = c.Expire(CONTROL_THREAD_LOCK, CONTROL_THREAD_EXPIRY)
			if err != nil {
				log.Printf("ERR: %s", err.Error())
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

func dropControlThread(c redis.Client) {
	val, err := c.Get(CONTROL_THREAD_LOCK)
	if err != nil {
		log.Printf("ERR: " + err.Error())
		return
	}
	if bytes.Compare(val, []byte(fmt.Sprint(config.Config.HostID))) == 0 {
		log.Printf("INFO: Found control thread with matching id %d, dropping", config.Config.HostID)
		_, err = c.Del(CONTROL_THREAD_LOCK)
		if err != nil {
			log.Printf("ERR: " + err.Error())
			return
		}
	}
}

func extendControlExpiry(c redis.Client) {
	c.Expire(CONTROL_THREAD_LOCK, CONTROL_THREAD_EXPIRY)
}

func threadControl() {
	if ControlThreadRunning {
		log.Printf("WARN: Control thread start attempting, but it looks like it's already running")
		return
	}

	log.Printf("INFO: Starting control thread")

	c, err := redis.NewSynchClientWithSpec(getConnection(REDIS_READWRITE).connspec)
	if err != nil {
		log.Printf("ERR: " + err.Error())
		return
	}

	// Check to see if we need to run control thread, or if it is currently
	// running on another host
	for {
		log.Printf("INFO: Attempting to grab control thread")
		if grabControlThread(c) {
			// If we're actually starting up, set global running flag
			ControlThreadRunning = true
			for {
				if !ControlThreadRunning {
					// Catch shutdown
					log.Printf("WARN: ControlThreadRunning was set to false, shutting down control thread")
					return
				}

				if util.ShuttingDown {
					log.Printf("INFO: ControlThread relinquishing control and terminating")
					dropControlThread(c)
					return
				}

				// Endlessly attempt to schedule checks
				log.Printf("DEBUG: attempt to schedule checks")
				members, err := c.Smembers(CHECKS_LIST)
				if err != nil {
					log.Printf("ERR: Unable to pull from key " + CHECKS_LIST)
				} else {
					// Pull list of all hosts/services
					for i := 0; i < len(members); i++ {
						// Pull last run and schedule interval to see if this needs to
						// be scheduled for another run, and push onto POLL_QUEUE.
						intervalRaw, err := c.Hget(string(members[i]), "interval")
						interval, _ := strconv.ParseUint(string(intervalRaw), 10, 64)

						command, err := c.Hget(string(members[i]), "command")

						typeRaw, err := c.Hget(string(members[i]), "type")
						checkType, _ := strconv.ParseUint(string(typeRaw), 10, 32)

						curtime := uint64(time.Now().Unix())

						items, err := c.Hgetall(string(members[i]) + ":hosts")
						if err == nil {
							for j := 0; j < len(items)/2; j += 2 {
								host := string(items[j])
								lastrun, _ := strconv.ParseUint(string(items[j+1]), 10, 64)
								//log.Info(fmt.Sprintf("curtime = %d, lastrun = %d, diff = %d, interval = %d", curtime, lastrun, curtime-lastrun, interval ))
								if curtime-lastrun >= interval {
									log.Printf("INFO: Adding %s : %s to poll queue", members[i], host)
									// Set lastrun to current time
									e := c.Hset(string(members[i])+":hosts", host, []byte(fmt.Sprint(curtime)))
									if e != nil {
										log.Printf("ERR: " + e.Error())
									}
									// Also update reverse index
									e = c.Hset(host+":checks", string(members[i]), []byte(fmt.Sprint(curtime)))
									if e != nil {
										log.Printf("ERR: " + e.Error())
									}
									// Form JSON object to serialize onto the scheduler stack
									obj := PollCheck{
										Host:        host,
										CheckName:   string(members[i]),
										EnqueueTime: curtime,
										Type:        uint(checkType),
										Command:     string(command),
									}
									o, err := json.Marshal(obj)
									if err == nil {
										e = c.Rpush(POLL_QUEUE, o)
										if e != nil {
											log.Printf("ERR: " + e.Error())
										}
									} else {
										log.Printf("ERR: " + e.Error())
									}
								}
							}
						}
					}

				}

				if util.ShuttingDown {
					log.Printf("INFO: ControlThread relinquishing control and terminating")
					dropControlThread(c)
					return
				}

				// Extend control thread expiry
				extendControlExpiry(c)

				// Sleep for a few seconds to avoid CPU piling.
				time.Sleep(2000 * time.Millisecond)
			}
		} else {
			log.Printf("DEBUG: ControlThread: already owned, waiting to start.")
		}

		// Sleep for a few seconds to avoid CPU piling.
		for i := 0; i < 15; i++ {
			if util.ShuttingDown {
				log.Printf("INFO: ControlThread relinquishing control and terminating")
				dropControlThread(c)
				return
			}
			time.Sleep(1 * time.Second)
		}
	} // for
}
