VERSION := $(shell git describe --tags --dirty --always)
LDFLAGS := -ldflags '-X "main.version=$(VERSION)" -s -w'

build:
	go build $(LDFLAGS) -o dist/mona
	cp README.md dist/README.md

test:
	go test -v ./...

lint:
	golangci-lint run