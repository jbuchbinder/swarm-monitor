// SWARM Distributed Monitoring System
// https://github.com/jbuchbinder/swarm-monitor
//
// vim: tabstop=2:softtabstop=2:shiftwidth=2:noexpandtab

package main

import (
	redis "github.com/jbuchbinder/go-redis"
	"fmt"
	"time"
)

func getContact(c redis.Client, name string) Contact {
	contact := Contact{}
	items, err := c.Hgetall(CONTACT_PREFIX + ":" + name)
	if err == nil {
		for j := 0; j < len(items)/2; j += 2 {
			k := string(items[j])
			v := string(items[j+1])
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
	c, cerr := redis.NewSynchClientWithSpec(getConnection(REDIS_READWRITE).connspec)
	if cerr != nil {
		log.Info(fmt.Sprintf("Alert thread #%d unable to acquire db connection", threadNum))
		return
	}
	log.Info(fmt.Sprintf("Starting alert thread #%d", threadNum))
	for {
		//log.Info(fmt.Sprintf("[%d] BLPOP %s 10", threadNum, ALERT_QUEUE))
		out, oerr := c.Blpop(ALERT_QUEUE, 0)
		if oerr != nil {
			log.Err(fmt.Sprintf("[ALERT %d] %s", threadNum, oerr.Error()))
		} else {
			if out == nil {
				log.Info(fmt.Sprintf("[ALERT %d] No output", threadNum))
			} else {
				if len(out) == 2 {
					log.Info(string(out[1]))
				}
			}
		}
		// Avoid potential pig-pile
		time.Sleep(10 * time.Millisecond)
	}
	return
}
