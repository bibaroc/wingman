regen:
	stringer -linecomment -type LogLevel ./pkg/logger

prepare:
	go get -u golang.org/x/perf/cmd/benchstat
	go get -u golang.org/x/tools/cmd/stringer
test:
	go test ./pkg/logger -benchmem -bench . - count 5 > tmp/logger.txt
	@benchstat tmp/logger.txt 