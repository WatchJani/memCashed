WORKING_DIR := ./cmd
TRACE_FLAG := --trace


run:
	@go run ./${WORKING_DIR}/main.go ${TRACE_FLAG}


test:
	@go test -bench BenchmarkSetReqPerSecond -benchtime 5s