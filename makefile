WORKING_DIR := ./cmd

run:
	


test_log:
	@go test ./cmd/server -bench BenchmarkSynchronousSet -benchmem &>> benchmarkSynchronous.log

test_duration:
	@go test ./cmd/server -bench BenchmarkSynchronousSet -benchtime 5s



benchmark_test:
	@go run ./${WORKING_DIR}/main.go > /dev/null 2>&1 &
	@go test ./cmd/server -bench BenchmarkSynchronousSet -benchmem && \
	go test ./cmd/server -bench BenchmarkSynchronousGet -benchmem && \
	go test ./cmd/server -bench BenchmarkSynchronousDelete -benchmem