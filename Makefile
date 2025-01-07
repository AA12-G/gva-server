.PHONY: dev
dev:
	fresh

.PHONY: build
build:
	go build -o ./tmp/main ./cmd/server

.PHONY: run
run:
	./tmp/main 