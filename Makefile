regen:
	stringer -linecomment -type LogLevel ./pkg/logger

prepare:
	go get -u golang.org/x/perf/cmd/benchstat
	go get -u golang.org/x/tools/cmd/stringer