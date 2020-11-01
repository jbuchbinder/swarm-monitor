// SWARM Distributed Monitoring System
// https://github.com/jbuchbinder/swarm-monitor
//
// vim: tabstop=2:softtabstop=2:shiftwidth=2:noexpandtab

package main

import (
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jbuchbinder/swarm-monitor/util"
)

func getContact(c redis.Client, name string) Contact {
	contact := Contact{}
	items, err := c.HGetAll(ctx, CONTACT_PREFIX+":"+name).Result()
	if err == nil {
		for k, v := range items {
			switch k {
			case "name":
				{
					contact.Name = v
				}
			case "email":
				{
					contact.EmailAddress = v
				}
			}
		}
	}
	return contact
}

func threadAlert(threadNum int) {
	c := redis.NewClient(getConnection(REDIS_READWRITE).connspec)

	log.Printf("INFO: Starting alert thread #%d", threadNum)

	util.RunningProcesses.Add(1)
	defer util.RunningProcesses.Done()

	for {
		if ControlThreadRunning {
			log.Printf("WARN: Control thread start attempting, but it looks like it's already running")
			return
		}

		log.Printf("DEBUG: [%d] BLPOP %s 10", threadNum, ALERT_QUEUE)
		out, oerr := c.BLPop(ctx, time.Duration(5)*time.Second, ALERT_QUEUE).Result()
		if oerr != nil {
			log.Printf("ERR: [ALERT %d] %s", threadNum, oerr.Error())
		} else {
			if out == nil {
				log.Printf("INFO: [ALERT %d] No output", threadNum)
			} else if len(out) == 2 {
				log.Printf("INFO: " + string(out[1]))
			}
		}
		// Avoid potential pig-pile
		time.Sleep(10 * time.Millisecond)
	}
}
