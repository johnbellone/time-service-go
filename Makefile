BUILD_DATE ?= $(shell date +%FT%T%z)
GIT_COMMIT ?= $(shell git rev-parse HEAD)
GIT_ABBRV ?= $(shell git describe --always --dirty --abbrev=8)
LDFLAGS = -X main.GitAbbrv=$(GIT_ABBRV) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_DATE)

.PHONY: all proto check build

all: build

clean:
	@rm -fr bin

check:
	@buf lint

proto: check
	@buf generate

build: proto
	@mkdir -p bin/
	@go build -ldflags "-s -w $(LDFLAGS)" -o bin/time-service

release: all
	@zip --junk-paths bin/* README.md LICENSE
