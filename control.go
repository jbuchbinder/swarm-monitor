package main

import (
	"./lib/redis"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

var (
	ControlThreadRunning = false
)

func threadControl() {
	if ControlThreadRunning {
		log.Warning("Control thread start attempting, but it looks like it's already running")
		return
	}

	log.Info("Starting control thread")

	c, err := redis.NewSynchClientWithSpec(getConnection(true).connspec)
	if err != nil {
		log.Err(err.Error())
		return
	}

	// Check to see if we need to run control thread, or if it is currently
	// running on another host

	// If we're actually starting up, set global running flag
	ControlThreadRunning = true
	for {
		if !ControlThreadRunning {
			// Catch shutdown
			log.Warning("ControlThreadRunning was set to false, shutting down control thread")
			return
		}

		// Endlessly attempt to schedule checks
		members, err := c.Smembers(CHECKS_LIST)
		if err != nil {
			log.Err("Unable to pull from key " + CHECKS_LIST)
		} else {
			// Pull list of all hosts/services
			for i := 0; i < len(members); i++ {
				// Pull last run and schedule interval to see if this needs to
				// be scheduled for another run, and push onto POLL_QUEUE.
				intervalRaw, err := c.Hget(string(members[i]), "interval")
				interval, _ := strconv.ParseUint(string(intervalRaw), 10, 64)
				curtime := uint64(time.Now().Unix())
				items, err := c.Hgetall(string(members[i]) + ":hosts")
				if err == nil {
					for j := 0; j < len(items)/2; j += 2 {
						host := string(items[j])
						lastrun, _ := strconv.ParseUint(string(items[j+1]), 10, 64)
						//log.Info(fmt.Sprintf("curtime = %d, lastrun = %d, diff = %d, interval = %d", curtime, lastrun, curtime-lastrun, interval ))
						if curtime-lastrun >= interval {
							log.Info(fmt.Sprintf("Adding %s : %s to poll queue", members[i], host))
							// Set lastrun to current time
							e := c.Hset(string(members[i])+":hosts", host, []byte(fmt.Sprint(curtime)))
							if e != nil {
								log.Err(e.Error())
							}
							// Form JSON object to serialize onto the scheduler stack
							obj := PollCheck{
								Host:        host,
								CheckName:   string(members[i]),
								EnqueueTime: curtime,
							}
							o, err := json.Marshal(obj)
							if err == nil {
								e = c.Rpush(POLL_QUEUE, o)
								if e != nil {
									log.Err(e.Error())
								}
							} else {
								log.Err(e.Error())
							}
						}
					}
				}
			}

		}

		// Sleep for a few seconds to avoid CPU piling.
		time.Sleep(2000 * time.Millisecond)
	}
}