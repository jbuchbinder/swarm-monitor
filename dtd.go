package main

import (
	"time"
)

const (
	CHECK_TYPE_NAGIOS = 1

	STATUS_OK       = 0
	STATUS_WARNING  = 1
	STATUS_CRITICAL = 2
	STATUS_UNKNOWN  = 3
)

type Status int

func (s Status) String() string {
	switch s {
	case STATUS_OK:
		return "OK"
	case STATUS_WARNING:
		return "WARNING"
	case STATUS_CRITICAL:
		return "CRITICAL"
	case STATUS_UNKNOWN:
		return "UNKNOWN"
	default:
		return "UNKNOWN"
	}
	return "UNKNOWN"
}

func NewStatus(s string) Status {
	switch s {
	case "OK":
		return STATUS_OK
	case "WARNING":
		return STATUS_WARNING
	case "CRITICAL":
		return STATUS_CRITICAL
	case "UNKNOWN":
		return STATUS_UNKNOWN
	default:
		return STATUS_UNKNOWN
	}
	return STATUS_UNKNOWN
}

type HostDefinition struct {
	Name    string   `json:"name"`
	Address string   `json:"address"`
	Checks  []string `json:"checks"`
}

type PollCheck struct {
	Host        string `json:"host"`
	Group       string `json:"group"`
	CheckName   string `json:"check"`
	Command     string `json:"command"`
	Type        uint   `json:"check_type"`
	EnqueueTime uint64 `json:"enqueue_time"`
}

type CheckStatus struct {
	Host              string    `json:"host"`
	CheckName         string    `json:"check"`
	CurrentStatus     Status    `json:"status"`
	CurrentStatusText string    `json:"status_text"`
	Since             time.Time `json:"since"`
	Iterations        int64     `json:"iterations"`
	LastOutput        string    `json:"output"`
}
