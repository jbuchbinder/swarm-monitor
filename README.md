SWARM-MONITOR
=============

Attempt to replace Nagios with a multi-threaded, clustered solution
based on Redis and Go.

REQUIREMENTS
------------

 * [Go](http://golang.org/) - compile the binary
 * [Redis](http://redis.io/) - run as db backend

BUILDING
--------

It's as simple as issuing a:

`go build`

LIBRARIES AND ACKNOWLEDGEMENTS
------------------------------

 * [go-redis](https://github.com/alphazero/Go-Redis) - AlphaZero's
   terrific Redis library for Go.
 * [gorest](http://code.google.com/p/gorest/) - REST library for Go.
   I ended up embedding this for simplicity, but you've got to give
   credit where credit is due.

