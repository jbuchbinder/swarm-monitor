package main

import (
	"./lib/redis"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func WaitTimeout(p *os.Process, timeout time.Duration) (*os.ProcessState, error) {
	timer := time.AfterFunc(timeout, func() { Kill(p) })
	defer timer.Stop()
	return p.Wait()
}

func Kill(p *os.Process) {
	syscall.Kill(p.Pid, syscall.SIGKILL)
}

func threadPoll(threadNum int) {
	log.Info(fmt.Sprintf("Starting poll thread #%d", threadNum))
	c, cerr := redis.NewSynchClientWithSpec(getConnection(true).connspec)
	if cerr != nil {
		log.Info(fmt.Sprintf("Poll thread #%d unable to acquire db connection", threadNum))
		return
	}
	for {
		//log.Info(fmt.Sprintf("[%d] BLPOP %s 10", threadNum, POLL_QUEUE))
		out, oerr := c.Blpop(POLL_QUEUE, 0)
		if oerr != nil {
			log.Err(fmt.Sprintf("[POLL %d] %s", threadNum, oerr.Error()))
		} else {
			if out == nil {
				log.Info(fmt.Sprintf("[ALERT %d] No output", threadNum))
			} else {
				if len(out) == 2 {
					check := PollCheck{}
					err := json.Unmarshal(out[1], &check)
					if err == nil {
						// Get check information

						// Process differently, depending on check type
						checkType := 1 // TODO: FIXME: FIXME: XXX: HACK
						switch {
						case checkType == CHECK_TYPE_NAGIOS:
							{
								// Do all appropriate substitutions
								cmd := &exec.Cmd{
									Path: "/usr/lib64/nagios/plugins/check_nrpe",
									Args: []string{"/usr/lib64/nagios/plugins/check_nrpe", "-H", "127.0.0.1", "-c", "check_disk"},
								}
								var bout bytes.Buffer
								cmd.Stdout = &bout
								err := cmd.Start()
								if err != nil {
									log.Err(err.Error())
								} else {
									// TODO: Configurable timeout for Nagios plugins
									msg, _ := WaitTimeout(cmd.Process, 30*time.Second)
									log.Info(fmt.Sprintf("Returned : %q\n", msg))
								}
							}
						}
					} else {
						log.Err(fmt.Sprintf("[ALERT %d] %s", threadNum, err.Error()))
					}
				}
			}
		}
		// Avoid potential pig-pile
		time.Sleep(10 * time.Millisecond)
	}
	return
}
