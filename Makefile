all: get-deps build

.PHONY: build
build:
	go build ./...

.PHONY: get-deps
get-deps:
	go get -v ./...

.PHONY: test
test: get-deps lint-check docs-check check

.PHONY: check
check:
	go test -v -race -cover ./...

.PHONY: lint-check
lint-check:
	golangci-lint run

.PHONY: fmt
fmt:
	@go fmt ./... | awk '{ print "Please run go fmt"; exit 1 }'

.PHONY: docs-dep
	which embedmd > /dev/null || go get github.com/campoy/embedmd

.PHONY: docs-check
docs-check: docs-dep
	@echo "Checking if docs are generated, if this fails, run 'make docs'."
	embedmd README.md | diff README.md -

.PHONY: docs
docs: docs-dep
	embedmd -w README.md
