package util

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	// ShutdownChannel listens for shutdown signals
	ShutdownChannel = make(chan os.Signal, 1)

	// ShuttingDown checks to see if application is terminating
	ShuttingDown = false
)

// SetCloseHandler starts up signal handlers for app termination
func SetCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal.Notify(ShutdownChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		ShuttingDown = true
		log.Printf("INFO: Ctrl+C pressed in Terminal")
		log.Printf("Waiting 10 seconds for threads to terminate")
		time.Sleep(time.Second * time.Duration(10))
		os.Exit(0)
	}()
	go log.Println(<-ShutdownChannel)
}

// WaitTimeout kills a process after a specified duration
func WaitTimeout(p *os.Process, timeout time.Duration) (*os.ProcessState, error) {
	timer := time.AfterFunc(timeout, func() { Kill(p) })
	defer timer.Stop()
	return p.Wait()
}

// Kill terminates a process
func Kill(p *os.Process) {
	syscall.Kill(p.Pid, syscall.SIGKILL)
}
