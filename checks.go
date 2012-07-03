// SWARM Distributed Monitoring System
// https://github.com/jbuchbinder/swarm-monitor
//
// vim: tabstop=2:softtabstop=2:shiftwidth=2:noexpandtab

package main

import (
	"./lib/redis"
)

func getMembersFromObj(c redis.Client, o string, p string) []string {
	entries, err := c.Smembers(o + ":" + p)
	if err != nil {
		log.Err(err.Error())
		ret := make([]string, 0)
		return ret
	}

	ret := make([]string, len(entries))
	for i := 0; i <= len(entries); i++ {
		ret[i] = string(entries[i])
	}
	return ret
}

func getAllMembersFromService(c redis.Client, o string) []string {
	hosts, err := c.Smembers(o + ":hosts")
	if err != nil {
		log.Err(err.Error())
	}
	h := make([]string, len(hosts))
	for i := 0; i <= len(hosts); i++ {
		h[i] = string(hosts[i])
	}
	groups, err := c.Smembers(o + ":groups")
	if err != nil {
		log.Err(err.Error())
	}
	for j := 0; j <= len(groups); j++ {
		m := getMembersFromObj(c, string(groups[j]), "members")
		for k := 0; k <= len(m); k++ {
			h = append(h, m[k])
		}
	}

	return h
}
