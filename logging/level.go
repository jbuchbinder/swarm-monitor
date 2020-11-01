package logging

import "strings"

const (
	LogTrace = iota
	LogDebug
	LogInfo
	LogWarn
	LogCrit
)

var (
	// LogLevel sets the current package logging level
	LogLevel = LogInfo
)

// StringToLevel converts a string into a logging level
func StringToLevel(s string) int {
	switch strings.ToUpper(s) {
	case "CRIT":
		return LogCrit
	case "CRITICAL":
		return LogCrit
	case "DEBUG":
		return LogDebug
	case "INFO":
		return LogInfo
	case "TRACE":
		return LogTrace
	case "WARN":
		return LogWarn
	case "WARNING":
		return LogWarn
	default:
		return LogWarn
	}
}
