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
	curl \
		--request POST http://localhost:8091/v1/yeet \
		--header "Content-Type: application/json" \
		--data '{"repository": "https://github.com/khuedoan/example-service", "revision": "master"}'
