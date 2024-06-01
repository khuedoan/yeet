.POSIX:
.PHONY: run temporal

run:
	go run .

temporal:
	temporal server start-dev
