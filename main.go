// SWARM Distributed Monitoring System
// https://github.com/jbuchbinder/swarm-monitor
//
// vim: tabstop=2:softtabstop=2:shiftwidth=2:noexpandtab

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jbuchbinder/swarm-monitor/config"
	"github.com/jbuchbinder/swarm-monitor/util"
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
	//log, _     = syslog.New(syslog.LOG_DEBUG, SERVICE_NAME)

	ctx context.Context
)

type redisConnection struct {
	host     string
	port     int
	password string
	db       int
	connspec *redis.Options
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

	c.connspec = &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.host, c.port),
		Password: c.password,
		DB:       c.db,
	}

	//redis.DefaultSpec().Host(c.host).Port(c.port).Db(c.db)
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

	ctx = context.Background()

	log.Printf("INFO: Starting " + SERVICE_NAME + " services")

	// Control thread
	log.Printf("INFO: Starting control thread")
	go threadControl()

	// Web thread
	log.Printf("INFO: Starting web thread")
	//go threadWeb()

	for t := 1; t <= config.Config.PoolSize; t++ {
		log.Printf("INFO: Attempting to spawn thread #%d", t)
		go threadAlert(t)
		go threadPoll(t)
	}

	util.SetCloseHandler()

	for {
		time.Sleep(10 * time.Second)
	}
}
