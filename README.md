# SWARM-MONITOR

[![Build Status](https://travis-ci.org/jbuchbinder/swarm-monitor.svg?branch=master)](https://travis-ci.org/jbuchbinder/swarm-monitor)
[![Go Report Card](https://goreportcard.com/badge/github.com/jbuchbinder/swarm-monitor)](https://goreportcard.com/report/github.com/jbuchbinder/swarm-monitor)
[![GoDoc](https://godoc.org/github.com/jbuchbinder/swarm-monitor?status.png)](https://godoc.org/github.com/jbuchbinder/swarm-monitor)

Attempt to replace Nagios with a multi-threaded, clustered solution
based on Redis and Go.

## REQUIREMENTS

 * [Go](http://golang.org/) - compile the binary
 * [Redis](http://redis.io/) - run as db backend

## BUILDING

It's as simple as issuing a:

`go build`

## LIBRARIES AND ACKNOWLEDGEMENTS

 * [go-redis](https://github.com/go-redis/redis) - Golang Redis library.
 * [gorest](http://code.google.com/p/gorest/) - REST library for Go. I ended up embedding this for simplicity, but you've got to give credit where credit is due.

