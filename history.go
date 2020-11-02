// SWARM Distributed Monitoring System
// https://github.com/jbuchbinder/swarm-monitor
//
// vim: tabstop=2:softtabstop=2:shiftwidth=2:noexpandtab

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type HistoryEvent struct {
	Timestamp   time.Time
	ID          int64
	HistoryKey  string
	SwarmHostID int
	Host        string
	Check       string
	Status      int32
	StatusText  string
	// TODO: FIXME: More keys
}

func (ev *HistoryEvent) PersistAtomicHistoryKV(c *redis.Client, k string, v []byte) {
	_, err := c.HSet(ctx, ev.HistoryKey, k, v).Result()
	if err != nil {
		log.Printf("ERR: " + fmt.Sprintf("Error persisting %s k %s v %s", ev.HistoryKey, k, v))
	}
}

func (ev *HistoryEvent) PersistEvent(c *redis.Client) {
	// Get new ever-incrementing event id
	id, err := c.Incr(ctx, HISTORY_KEY).Result()
	if err != nil {
		log.Printf("ERR: " + err.Error())
		return
	}

	// Persist to object internally
	ev.ID = id

	// New key
	historyKey := HISTORY_BASE + ":id:" + fmt.Sprintf("%d", id)
	ev.HistoryKey = historyKey

	// Persist values to history key
	ev.PersistAtomicHistoryKV(c, "timestamp", []byte(ev.Timestamp.String()))
	// TODO: FIXME: More keys

	// Build additional indices...
	// 1. Master index.
	_, err = c.ZAdd(ctx, HISTORY_LIST, &redis.Z{
		Score:  float64(ev.Timestamp.Unix()),
		Member: []byte(historyKey),
	}).Result()
	if err != nil {
		log.Printf("ERR: " + err.Error())
	}

	// 2. Date index
	log.Printf("INFO: Logging to " + HISTORY_LIST + ":date:" + ev.Timestamp.Format("2006-01-02"))
	_, err = c.ZAdd(ctx, HISTORY_LIST+":date:"+ev.Timestamp.Format("2006-01-02"), &redis.Z{
		Score:  float64(ev.Timestamp.Unix()),
		Member: []byte(historyKey),
	}).Result()
	if err != nil {
		log.Printf("ERR: " + err.Error())
	}

	// 3. Host index
	log.Printf("INFO: Logging to " + HISTORY_LIST + ":host:" + ev.Host)
	_, err = c.ZAdd(ctx, HISTORY_LIST+":host:"+ev.Host, &redis.Z{
		Score:  float64(ev.Timestamp.Unix()),
		Member: []byte(historyKey),
	}).Result()
	if err != nil {
		log.Printf("ERR: " + err.Error())
	}

	// 4. Check index
	log.Printf("INFO: Logging to " + HISTORY_LIST + ":check:" + ev.Check)
	_, err = c.ZAdd(ctx, HISTORY_LIST+":check:"+ev.Check, &redis.Z{
		Score:  float64(ev.Timestamp.Unix()),
		Member: []byte(historyKey),
	}).Result()
	if err != nil {
		log.Printf("ERR: " + err.Error())
	}

}
