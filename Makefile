# If using multiple versions of go, or your target golang version is not in your
# PATH, or whatever reason, specify a path to the golang binary when invoking
# make, for example: make build GO=/path/to/other/golang/bin/go
GO ?= go

all: deps build

.PHONY: build
build:
	$(GO) build ./...

.PHONY: deps
deps:
	$(GO) mod tidy && $(GO) mod download && $(GO) mod verify

.PHONY: test
test: deps check

.PHONY: check
check:
	$(GO) test -v -race -cover ./...

.PHONY: fmt
fmt:
	@$(GO) fmt ./... | awk '{ print "Please run go fmt"; exit 1 }'
