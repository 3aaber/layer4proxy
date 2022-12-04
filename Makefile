GO=go
BIN_DIR:=bin
SOURCES=$(shell find . -name '*.go' -not -name '*_test.go' -not -name "main.go")

all: pre $(BIN_DIR)/layer4proxy

layer4proxy: pre $(BIN_DIR)/layer4proxy

pre:
	mkdir -p $(BIN_DIR)

clean:
	rm -rf $(BIN_DIR)

format:
	find . -name '*.go' -not -path './vendor/*' | xargs -n1 go fmt

$(BIN_DIR)/%: cmd/% $(SOURCES)
	$(GO) build -ldflags="-s -w" -o $@ $</*.go

$(BIN_DIR):
	mkdir -p $@

test:
	$(GO) test -v ./...

coverage.out: $(wildcard **/**/*.go)
	$(GO) test -cover -coverprofile=coverage.out ./...

coverage.html: coverage.out
	$(GO) tool cover -html=coverage.out -o coverage.html

dev-redis:
	docker-compose -f build/redis-sentinel/docker-compose.yml up -d

.PHONY: all master client pre clean format test dev-redis
