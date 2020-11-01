package util

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	// ShutdownChannel listens for shutdown signals
	ShutdownChannel = make(chan os.Signal, 1)

	// ShuttingDown checks to see if application is terminating
	ShuttingDown = false

	// RunningProcesses describes a wait group
	RunningProcesses = &sync.WaitGroup{}
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
		log.Printf("Waiting for threads to terminate")
		RunningProcesses.Wait()
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
