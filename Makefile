regen:
	stringer -linecomment -type LogLevel ./pkg/logger

prepare:
	go get -u golang.org/x/perf/cmd/benchstat
	go get -u golang.org/x/tools/cmd/stringer
	go get -u github.com/google/go-cmp/cmp

bench:
	go test ./pkg/logger -benchmem -bench . -run '' -count 5 | tee tmp/logger.txt
	@benchstat tmp/logger.txt
	go test ./pkg/fst -benchmem -bench . -run '' -count 5 | tee tmp/fst.txt
	@benchstat tmp/fst.txt