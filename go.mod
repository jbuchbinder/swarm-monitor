module github.com/jbuchbinder/swarm-monitor

go 1.15

replace (
	github.com/jbuchbinder/swarm-monitor/checks => ./checks
	github.com/jbuchbinder/swarm-monitor/config => ./config
	github.com/jbuchbinder/swarm-monitor/logging => ./logging
	github.com/jbuchbinder/swarm-monitor/util => ./util
)

require (
	github.com/ajg/form v1.5.1 // indirect
	github.com/fromkeith/gorest v0.0.0-20150514202557-ee389d6398d5
	github.com/go-redis/redis/v8 v8.3.3
	github.com/gorilla/schema v1.2.0 // indirect
	github.com/jbuchbinder/go-redis v0.0.0-20130426131933-cdc4c1e04f41
	github.com/jbuchbinder/swarm-monitor/checks v0.0.0-00010101000000-000000000000
	github.com/jbuchbinder/swarm-monitor/config v0.0.0-00010101000000-000000000000
	github.com/jbuchbinder/swarm-monitor/logging v0.0.0-00010101000000-000000000000
	github.com/jbuchbinder/swarm-monitor/util v0.0.0-00010101000000-000000000000
)
