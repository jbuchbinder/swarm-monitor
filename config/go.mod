module github.com/jbuchbinder/swarm-monitor/config

go 1.16

replace github.com/jbuchbinder/swarm-monitor/logging => ../logging

require (
	github.com/jbuchbinder/swarm-monitor/logging v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v2 v2.3.0
)
