package logging

import (
	"fmt"
	"log"
	"runtime"
)

// Critf prints CRIT level formatting
func Critf(format string, args ...interface{}) {
	if LogLevel <= LogCrit {
		log.Printf("CRIT: "+getCaller()+": "+format, args...)
	}
}

// Debugf prints DEBUG level formatting
func Debugf(format string, args ...interface{}) {
	if LogLevel <= LogDebug {
		log.Printf("DEBUG: "+getCaller()+": "+format, args...)
	}
}

// Infof prints INFO level formatting
func Infof(format string, args ...interface{}) {
	if LogLevel <= LogInfo {
		log.Printf("INFO: "+getCaller()+": "+format, args...)
	}
}

// Tracef prints TRACE level formatting
func Tracef(format string, args ...interface{}) {
	if LogLevel <= LogTrace {
		log.Printf("TRACE: "+getCaller()+": "+format, args...)
	}
}

// Warnf prints WARN level formatting
func Warnf(format string, args ...interface{}) {
	if LogLevel <= LogWarn {
		log.Printf("WARN: "+getCaller()+": "+format, args...)
	}
}

func getCaller() string {
	_, file, no, _ := runtime.Caller(2)
	return fmt.Sprintf("%s:%d", file, no)
}
