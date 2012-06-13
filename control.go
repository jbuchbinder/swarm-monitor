package main

import (
	"./lib/redis"
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
			}

		}

		// Sleep for a few seconds to avoid CPU piling.
		time.Sleep(300 * time.Millisecond)
	}
}
