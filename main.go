package main

import (
	"./lib/redis"
	"flag"
	"fmt"
	"log/syslog"
	"time"
)

const (
	SERVICE_NAME          = "swarm"
	CONTROL_THREAD_LOCK   = SERVICE_NAME + ":lock:control"
	ALERT_QUEUE           = SERVICE_NAME + ":queue:alert"
	POLL_QUEUE            = SERVICE_NAME + ":queue:poll"
	CHECKS_LIST           = SERVICE_NAME + ":index:checks"
	CONTROL_THREAD_EXPIRY = 60
)

var (
	redisHost = flag.String("dbhost", "localhost", "Redis host")
	redisPort = flag.Int("dbport", 6379, "Redis port")
	redisDb   = flag.Int("dbnum", 13, "Redis db number")
	poolSize  = flag.Int("pool", 5, "Thread pool size")
	webPort   = flag.Int("webport", 48666, "Web listening port")
	hostId    = flag.Int("hostid", 1, "Server host id for cluster")
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
