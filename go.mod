module github.com/jbuchbinder/swarm-monitor

go 1.15

replace github.com/jbuchbinder/swarm-monitor/config => ./config

require (
	github.com/ajg/form v1.5.1 // indirect
	github.com/fromkeith/gorest v0.0.0-20150514202557-ee389d6398d5
	github.com/gorilla/schema v1.2.0 // indirect
	github.com/jbuchbinder/go-redis v0.0.0-20130426131933-cdc4c1e04f41
	github.com/jbuchbinder/swarm-monitor/config v0.0.0-00010101000000-000000000000
)
