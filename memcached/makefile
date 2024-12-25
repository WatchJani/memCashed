test_log:
	@go test ./server -bench BenchmarkSynchronousSet -benchmem &>> benchmarkSynchronous.log

test_duration:
	@go test ./server -bench BenchmarkSynchronousSet -benchtime 5s

benchmark_test:
	@go run ./main.go > /dev/null 2>&1 &
	@go test ./server -bench BenchmarkSynchronousSet -benchmem && \
	go test ./server -bench BenchmarkSynchronousGet -benchmem && \
	go test ./server -bench BenchmarkSynchronousDelete -benchmem

build:
	go build -o memcached