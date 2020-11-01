// SWARM Distributed Monitoring System
// https://github.com/jbuchbinder/swarm-monitor
//
// vim: tabstop=2:softtabstop=2:shiftwidth=2:noexpandtab

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	//redis "github.com/jbuchbinder/go-redis"
	//redis "github.com/funkygao/Go-Redis"
	"github.com/jbuchbinder/swarm-monitor/checks"
	"github.com/jbuchbinder/swarm-monitor/util"
)

func updateCheckResults(c *redis.Client, host string, check string, status int32, statusText string) {
	log.Printf("INFO: updateCheckResults %s %s : %d [%s]", host, check, status, statusText)
	// TODO: persist update check results
}

func threadPoll(threadNum int) {
	log.Printf("INFO: Starting poll thread #%d", threadNum)

	c := redis.NewClient(getConnection(REDIS_READWRITE).connspec)

	util.RunningProcesses.Add(1)
	defer util.RunningProcesses.Done()

	for {
		if util.ShuttingDown {
			log.Printf("INFO: thread %d shutting down", threadNum)
			return
		}

		log.Printf("INFO: [%d] BLPOP %s 10", threadNum, POLL_QUEUE)
		out, oerr := c.BLPop(ctx, time.Duration(5)*time.Second, POLL_QUEUE).Result()
		if oerr != nil {
			log.Printf("ERR: [POLL %d] %s", threadNum, oerr.Error())
		} else {
			if out == nil {
				log.Printf("INFO: [ALERT %d] No output", threadNum)
			} else {
				if len(out) == 2 {
					check := PollCheck{}
					err := json.Unmarshal([]byte(out[1]), &check)
					if err == nil {
						// Get check information

						// Process differently, depending on check type
						checkType := check.Type
						switch {
						case checkType == checks.CheckTypeBuiltIn:
							{
								chk, err := checks.InstantiateChecker(check.Command)
								if err != nil {
									log.Printf("ERR: " + err.Error())
									break
								}
								exitStatus, msg := chk.Check(check.Host)
								log.Printf("INFO: " + fmt.Sprintf("Returned : %d:%q\n", exitStatus, msg))
								updateCheckResults(c, check.Host, check.CheckName, int32(exitStatus), msg)
							}
						case checkType == checks.CheckTypeNagios:
							{
								// TODO: Handle additional options, substitutions and overrides
								replacer := util.ReplacerFromMap(map[string]string{
									"$HOSTADDRESS$": check.Host,
								})
								cmdParts := strings.Split(replacer.Replace(check.Command), " ")

								// Do all appropriate substitutions
								cmd := &exec.Cmd{
									Path: cmdParts[0],
									Args: cmdParts,
								}
								var bout bytes.Buffer
								cmd.Stdout = &bout
								err := cmd.Start()
								if err != nil {
									log.Printf("ERR: " + err.Error())
								} else {
									// TODO: Configurable timeout for Nagios plugins
									var exitStatus int
									msg, cerr := util.WaitTimeout(cmd.Process, 30*time.Second)
									if cerr != nil {
										// Handle timeout

										// WARNING: This is UNIX/Linux specific. For Windows,
										// this would need to be followed:
										// https://groups.google.com/d/msg/golang-nuts/8XIlxWgpdJw/Z8s2N-SoWHsJ

										if msg, ok := err.(*exec.ExitError); ok { // there is error code
											exitStatus = msg.Sys().(syscall.WaitStatus).ExitStatus()
										} else {
											// Okay status
											exitStatus = 0
										}
									} else {
										// Handle return status
										exitStatus = 0
									}
									log.Printf("INFO: Returned : %d:%q\n", exitStatus, msg)
									updateCheckResults(c, check.Host, check.CheckName, int32(exitStatus), msg.String())
								}
							}
						}
					} else {
						log.Printf("ERR: [ALERT %d] %s", threadNum, err.Error())
					}
				}
			}
		}
		// Avoid potential pig-pile
		time.Sleep(10 * time.Millisecond)
	}
}
