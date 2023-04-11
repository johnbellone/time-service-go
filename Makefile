BUILD_DATE ?= $(shell date +%FT%T%z)
GIT_COMMIT ?= $(shell git rev-parse HEAD)
GIT_ABBRV ?= $(shell git describe --always --dirty --abbrev=8)
LDFLAGS = -X main.GitAbbrv=$(GIT_ABBRV) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_DATE)

.PHONY: all proto check build tools

all: build

clean:
	@rm -fr bin

tools:
	@go mod tidy
	@go generate -tags tools tools/tools.go
	@minica -domains localhost -ip-addresses 127.0.0.1 -ca-cert server.crt -ca-key server.key
	@openssl rsa -in server.key -pubout > server.pub 2>&1
	@rm -fr localhost

check:
	@buf lint

proto: check
	@buf generate

build:
	@mkdir -p bin/
	@go build -ldflags "$(LDFLAGS)" -o bin/time-service

release: all
	@zip --junk-paths bin/* README.md LICENSE
