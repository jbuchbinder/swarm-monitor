package main

import (
	"./lib/redis"
	"flag"
	"fmt"
	"log/syslog"
	"time"
)

const (
	SERVICE_NAME = "monitor"
	ALERT_QUEUE  = SERVICE_NAME + ":queue:alert"
	POLL_QUEUE   = SERVICE_NAME + ":queue:poll"
	CHECKS_LIST  = SERVICE_NAME + ":index:checks"
)

var (
	redisHost = flag.String("dbhost", "localhost", "Redis host")
	redisPort = flag.Int("dbport", 6379, "Redis port")
	redisDb   = flag.Int("dbnum", 13, "Redis db number")
	poolSize  = flag.Int("pool", 5, "Thread pool size")
	webPort   = flag.Int("webport", 48666, "Web listening port")
	log, _    = syslog.New(syslog.LOG_DEBUG, SERVICE_NAME)
)

type redisConnection struct {
	host     string
	port     int
	password string
	db       int
	connspec *redis.ConnectionSpec
}

func getConnection(write bool) redisConnection {
	var c redisConnection

	if write {
		c.host = *redisHost
		c.port = *redisPort
		c.db = *redisDb
		//c.password = ""
	} else {
		c.host = *redisHost
		c.port = *redisPort
		c.db = *redisDb
		//c.password = ""
	}

	c.connspec = redis.DefaultSpec().Host(c.host).Port(c.port).Db(c.db)
	//.Password(c.password)

	return c
}

func threadAlert(threadNum int) {
	c, cerr := redis.NewSynchClientWithSpec(getConnection(true).connspec)
	if cerr != nil {
		log.Info(fmt.Sprintf("Poll thread #%d unable to acquire db connection", threadNum))
		return
	}
	log.Info(fmt.Sprintf("Starting alert thread #%d", threadNum))
	for {
		log.Info(fmt.Sprintf("[%d] BLPOP %s 10", threadNum, ALERT_QUEUE))
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

func threadPoll(threadNum int) {
	log.Info(fmt.Sprintf("Starting poll thread #%d", threadNum))
	c, cerr := redis.NewSynchClientWithSpec(getConnection(true).connspec)
	if cerr != nil {
		log.Info(fmt.Sprintf("Poll thread #%d unable to acquire db connection", threadNum))
		return
	}
	for {
		log.Info(fmt.Sprintf("[%d] BLPOP %s 10", threadNum, POLL_QUEUE))
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

func main() {
	flag.Parse()
	log.Info("Starting " + SERVICE_NAME + " services")

	// Control thread
	go threadControl()

	// Web thread
	go threadWeb()

	for t := 1; t <= *poolSize; t++ {
		log.Info(fmt.Sprintf("Attempting to spawn thread #%d", t))
		go threadAlert(t)
		go threadPoll(t)
	}

	for {
		time.Sleep(10 * time.Second)
	}
}
