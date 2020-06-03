regen:
	stringer -linecomment -type LogLevel ./pkg/logger

prepare:
	go get -u golang.org/x/perf/cmd/benchstat
	go get -u golang.org/x/tools/cmd/stringer

bench:
	go test ./pkg/logger -benchmem -bench . -run '' -count 5 | tee tmp/logger.txt
	@benchstat tmp/logger.txt