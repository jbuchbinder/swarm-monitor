package main

import (
	"./lib/redis"
	"fmt"
	"time"
)

func threadAlert(threadNum int) {
	c, cerr := redis.NewSynchClientWithSpec(getConnection(true).connspec)
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
