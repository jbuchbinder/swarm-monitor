// SWARM Distributed Monitoring System
// https://github.com/jbuchbinder/swarm-monitor
//
// vim: tabstop=2:softtabstop=2:shiftwidth=2:noexpandtab

package main

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jbuchbinder/swarm-monitor/logging"
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

	logging.Infof("Starting alert thread #%d", threadNum)

	util.RunningProcesses.Add(1)
	defer util.RunningProcesses.Done()

	for {
		if ControlThreadRunning {
			logging.Debugf("Control thread start attempting, but it looks like it's already running")
			return
		}

		logging.Debugf("[%d] BLPOP %s 10", threadNum, ALERT_QUEUE)
		out, oerr := c.BLPop(ctx, time.Duration(5)*time.Second, ALERT_QUEUE).Result()
		if oerr != nil {
			logging.Debugf("[ALERT %d] %s", threadNum, oerr.Error())
		} else {
			if out == nil {
				logging.Infof("[ALERT %d] No output", threadNum)
			} else if len(out) == 2 {
				logging.Infof(string(out[1]))
			}
		}
		// Avoid potential pig-pile
		time.Sleep(10 * time.Millisecond)
	}
}
