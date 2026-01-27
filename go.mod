module github.com/jbuchbinder/swarm-monitor

go 1.24.0

replace (
	github.com/jbuchbinder/swarm-monitor/checks => ./checks
	github.com/jbuchbinder/swarm-monitor/config => ./config
	github.com/jbuchbinder/swarm-monitor/logging => ./logging
	github.com/jbuchbinder/swarm-monitor/util => ./util
)

require (
	github.com/fromkeith/gorest v0.0.0-20150514202557-ee389d6398d5
	github.com/go-redis/redis/v8 v8.11.5
	github.com/jbuchbinder/go-redis v0.0.0-20130426131933-cdc4c1e04f41
	github.com/jbuchbinder/swarm-monitor/checks v0.0.0-20210302005737-448cf807518b
	github.com/jbuchbinder/swarm-monitor/config v0.0.0-20210302005737-448cf807518b
	github.com/jbuchbinder/swarm-monitor/logging v0.0.0-20210302005737-448cf807518b
	github.com/jbuchbinder/swarm-monitor/util v0.0.0-20210302005737-448cf807518b
)

require (
	github.com/ajg/form v1.5.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/chzyer/logex v1.1.10 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/chzyer/test v0.0.0-20180213035817-a1ea475d72b1 // indirect
	github.com/creack/pty v1.1.9 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/pprof v0.0.0-20210407192527-94a9f03dee38 // indirect
	github.com/gorilla/schema v1.4.1 // indirect
	github.com/hpcloud/tail v1.0.0 // indirect
	github.com/ianlancetaylor/demangle v0.0.0-20200824232613-28f6c0f3b639 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/pty v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/ginkgo/v2 v2.0.0 // indirect
	github.com/onsi/gomega v1.18.1 // indirect
	github.com/pkg/diff v0.0.0-20210226163009-20ebb0f2a09e // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	github.com/yuin/goldmark v1.4.13 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.39.0 // indirect
	go.opentelemetry.io/otel/metric v1.39.0 // indirect
	go.opentelemetry.io/otel/trace v1.39.0 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/mod v0.21.0 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/telemetry v0.0.0-20240521205824-bda55230c457 // indirect
	golang.org/x/term v0.25.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	golang.org/x/tools v0.26.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
