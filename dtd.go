// SWARM Distributed Monitoring System
// https://github.com/jbuchbinder/swarm-monitor
//
// vim: tabstop=2:softtabstop=2:shiftwidth=2:noexpandtab

package main

import (
	"time"

	"github.com/jbuchbinder/swarm-monitor/checks"
)

type Status int

func (s Status) String() string {
	switch s {
	case checks.StatusOK:
		return "OK"
	case checks.StatusWarning:
		return "WARNING"
	case checks.StatusCritical:
		return "CRITICAL"
	case checks.StatusUnknown:
		return "UNKNOWN"
	default:
		return "UNKNOWN"
	}
}

func NewStatus(s string) Status {
	switch s {
	case "OK":
		return checks.StatusOK
	case "WARNING":
		return checks.StatusWarning
	case "CRITICAL":
		return checks.StatusCritical
	case "UNKNOWN":
		return checks.StatusUnknown
	default:
		return checks.StatusUnknown
	}
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

type Contact struct {
	Name         string `json:"name"`
	DisplayName  string `json:"display_name"`
	EmailAddress string `json:"email"`
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
