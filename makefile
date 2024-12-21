WORKING_DIR := ./cmd
TRACE_FLAG := --trace


run:
	@go run ./${WORKING_DIR}/main.go ${TRACE_FLAG}


test_log:
	@go test ./cmd/server -bench BenchmarkSynchronous -benchmem &>> benchmarkConcatenation.log

test_duration:
	@go test ./cmd/server -bench BenchmarkSynchronous -benchtime 5s