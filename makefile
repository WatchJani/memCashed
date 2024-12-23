WORKING_DIR := ./cmd
TRACE_FLAG := --trace


run:
	@go run ./${WORKING_DIR}/main.go ${TRACE_FLAG}


test_log:
	@go test ./cmd/server -bench BenchmarkSynchronousSet -benchmem &>> benchmarkSynchronous.log

test_duration:
	@go test ./cmd/server -bench BenchmarkSynchronousSet -benchtime 5s