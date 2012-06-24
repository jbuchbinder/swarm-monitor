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
	err = c.Rpush(HISTORY_LIST, []byte(historyKey))
	if err != nil {
		log.Err(err.Error())
	}

	// 2. Date index
	log.Info("Logging to " + HISTORY_LIST + ":date:" + ev.Timestamp.Format("%Y-%m-%d"))
	err = c.Rpush(HISTORY_LIST+":date:"+ev.Timestamp.Format("%Y-%m-%d"), []byte(historyKey))
	if err != nil {
		log.Err(err.Error())
	}

	// 3. Host index
	// TODO: FIXME

	// 4. Service index
	// TODO: FIXME
}
