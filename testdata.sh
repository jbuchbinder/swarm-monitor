#!/bin/bash

TMP=$(mktemp XXXXXXXXXX)

cat<<EOF > $TMP
ECHO "Clear database"
SELECT 13
FLUSHDB
ECHO "Populate hosts"
HMSET monitor:hosts:hmon01 name hmon01
SADD monitor:index:hosts monitor:hosts:hmon01
HMSET monitor:hosts:hmon02 name hmon02
SADD monitor:index:hosts monitor:hosts:hmon02
ECHO "Populate checks"
HMSET monitor:checks:check_disk name check_disk interval 300
HMSET monitor:checks:check_disk:hosts monitor:hosts:hmon01 0 monitor:hosts:hmon02 0
SADD monitor:index:checks monitor:checks:check_disk
EOF

cat $TMP | /usr/bin/redis-cli

rm -f $TMP

