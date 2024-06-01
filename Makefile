.POSIX:
.PHONY: temporal server worker fmt

temporal:
	temporal server start-dev

server:
	go run ./server

worker:
	go run ./worker

fmt:
	go fmt ./...

test:
	go test ./...
	curl -X POST http://localhost:8091/v1/yeet
