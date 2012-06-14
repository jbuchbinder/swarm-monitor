package main

import (
	"./lib/redis"
	"fmt"
	"time"
)

func threadPoll(threadNum int) {
	log.Info(fmt.Sprintf("Starting poll thread #%d", threadNum))
	c, cerr := redis.NewSynchClientWithSpec(getConnection(true).connspec)
	if cerr != nil {
		log.Info(fmt.Sprintf("Poll thread #%d unable to acquire db connection", threadNum))
		return
	}
	for {
		//log.Info(fmt.Sprintf("[%d] BLPOP %s 10", threadNum, POLL_QUEUE))
		out, oerr := c.Blpop(POLL_QUEUE, 0)
		if oerr != nil {
			log.Err(fmt.Sprintf("[POLL %d] %s", threadNum, oerr.Error()))
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
