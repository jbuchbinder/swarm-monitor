#!/bin/bash

TMP=$(mktemp XXXXXXXXXX)

cat<<'EOF' > $TMP
ECHO "Clear database"
SELECT 13
FLUSHDB
ECHO "Populate hosts"
HMSET swarm:hosts:hmon01 name hmon01
HMSET swarm:hosts:hmon01 address hmon01
SADD swarm:index:hosts swarm:hosts:hmon01
HMSET swarm:hosts:hmon02 name hmon02
HMSET swarm:hosts:hmon02 address hmon02
SADD swarm:index:hosts swarm:hosts:hmon02
ECHO "Populate checks"
HMSET swarm:checks:check_disk name check_disk interval 300 timeout 30
HMSET swarm:checks:check_disk type 1 command "/usr/lib64/nagios/plugins/check_nrpe -H $HOSTADDRESS$ -c check_disk -t 30"
HMSET swarm:checks:check_disk:hosts swarm:hosts:hmon01 0 swarm:hosts:hmon02 0
HMSET swarm:hosts:hmon01:checks swarm:checks:check_disk 0
HMSET swarm:hosts:hmon02:checks swarm:checks:check_disk 0
SADD swarm:index:checks swarm:checks:check_disk
ECHO "Populate contacts"
SADD swarm:index:contacts swarm:contacts:jbuchbinder
HMSET swarm:contacts:jbuchbinder name jbuchbinder display_name "Jeff Buchbinder" email "jeff@ourexchange.net"
EOF

cat $TMP | /usr/bin/redis-cli

rm -f $TMP

