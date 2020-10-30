// SWARM Distributed Monitoring System
// https://github.com/jbuchbinder/swarm-monitor
//
// vim: tabstop=2:softtabstop=2:shiftwidth=2:noexpandtab

package main

import (
	"flag"
	"fmt"
	redis "github.com/jbuchbinder/go-redis"
	"github.com/jbuchbinder/swarm-monitor/config"
	"log/syslog"
	"time"
)

const (
	SERVICE_NAME          = "swarm"
	CONTROL_THREAD_LOCK   = SERVICE_NAME + ":lock:control"
	ALERT_QUEUE           = SERVICE_NAME + ":queue:alert"
	POLL_QUEUE            = SERVICE_NAME + ":queue:poll"
	HOSTS_LIST            = SERVICE_NAME + ":index:hosts"
	HOST_PREFIX           = SERVICE_NAME + ":hosts"
	CHECK_PREFIX          = SERVICE_NAME + ":checks"
	CHECKS_LIST           = SERVICE_NAME + ":index:checks"
	HISTORY_BASE          = SERVICE_NAME + ":history"
	HISTORY_LIST          = HISTORY_BASE + ":index"
	HISTORY_KEY           = HISTORY_BASE + ":key"
	CONTACT_PREFIX        = SERVICE_NAME + ":contact"
	CONTACT_LIST          = SERVICE_NAME + ":index:contacts"
	CONTROL_THREAD_EXPIRY = 60
	REDIS_READONLY        = false
	REDIS_READWRITE       = true
)

var (
	configFile = flag.String("config", "swarm.yml", "YAML configuration file")
	log, _     = syslog.New(syslog.LOG_DEBUG, SERVICE_NAME)
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
		c.host = config.Config.Database.Host
		c.port = config.Config.Database.Port
		c.db = config.Config.Database.Db
		//c.password = ""
	} else {
		c.host = config.Config.Database.Host
		c.port = config.Config.Database.Port
		c.db = config.Config.Database.Db
		//c.password = ""
	}

	c.connspec = redis.DefaultSpec().Host(c.host).Port(c.port).Db(c.db)
	//.Password(c.password)

	return c
}

func main() {
	flag.Parse()

	c, err := config.LoadConfigWithDefaults(*configFile)
	if err != nil {
		panic(err)
	}
	if c == nil {
		panic("UNABLE TO LOAD CONFIG")
	}
	config.Config = c

	log.Info("Starting " + SERVICE_NAME + " services")

	// Control thread
	go threadControl()

	// Web thread
	go threadWeb()

	for t := 1; t <= config.Config.PoolSize; t++ {
		log.Info(fmt.Sprintf("Attempting to spawn thread #%d", t))
		go threadAlert(t)
		go threadPoll(t)
	}

	for {
		time.Sleep(10 * time.Second)
	}
}
