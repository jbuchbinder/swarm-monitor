package main

const (
	CHECK_TYPE_NAGIOS = 1
)

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
