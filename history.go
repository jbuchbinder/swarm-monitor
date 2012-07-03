// SWARM Distributed Monitoring System
// https://github.com/jbuchbinder/swarm-monitor
//
// vim: tabstop=2:softtabstop=2:shiftwidth=2:noexpandtab

package main

import (
	"./lib/redis"
	"fmt"
	"time"
)

type HistoryEvent struct {
	Timestamp  time.Time
	Id         int64
	HistoryKey string
	Host       string
	Check      string
	// TODO: FIXME: More keys
}

func (ev *HistoryEvent) PersistAtomicHistoryKV(c redis.Client, k string, v []byte) {
	err := c.Hset(ev.HistoryKey, k, v)
	if err != nil {
		log.Err(fmt.Sprintf("Error persisting %s k %s v %s", ev.HistoryKey, k, v))
	}
}

func (ev *HistoryEvent) persistEvent(c redis.Client) {
	// Get new ever-incrementing event id
	id, err := c.Incr(HISTORY_KEY)
	if err != nil {
		log.Err(err.Error())
		return
	}

	// Persist to object internally
	ev.Id = id

	// New key
	historyKey := HISTORY_BASE + ":id:" + string(id)
	ev.HistoryKey = historyKey

	// Persist values to history key
	ev.PersistAtomicHistoryKV(c, "timestamp", []byte(ev.Timestamp.String()))
	// TODO: FIXME: More keys

	// Build additional indices...
	// 1. Master index.
	_, err = c.Zadd(HISTORY_LIST, float64(ev.Timestamp.Unix()), []byte(historyKey))
	if err != nil {
		log.Err(err.Error())
	}

	// 2. Date index
	log.Info("Logging to " + HISTORY_LIST + ":date:" + ev.Timestamp.Format("%Y-%m-%d"))
	_, err = c.Zadd(HISTORY_LIST+":date:"+ev.Timestamp.Format("%Y-%m-%d"), float64(ev.Timestamp.Unix()), []byte(historyKey))
	if err != nil {
		log.Err(err.Error())
	}

	// 3. Host index
	log.Info("Logging to " + HISTORY_LIST + ":host:" + ev.Host)
	_, err = c.Zadd(HISTORY_LIST+":host:"+ev.Host, float64(ev.Timestamp.Unix()), []byte(historyKey))
	if err != nil {
		log.Err(err.Error())
	}

	// 4. Check index
	log.Info("Logging to " + HISTORY_LIST + ":check:" + ev.Check)
	_, err = c.Zadd(HISTORY_LIST+":check:"+ev.Check, float64(ev.Timestamp.Unix()), []byte(historyKey))
	if err != nil {
		log.Err(err.Error())
	}

}
