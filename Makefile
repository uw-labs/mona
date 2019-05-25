VERSION := $(shell git describe --tags --dirty --always)
LDFLAGS := -ldflags '-X "main.version=$(VERSION)" -s -w'
IMAGE_NAME := davidsbond/mona

build:
	go build $(LDFLAGS) -o dist/mona
	cp README.md dist/README.md
	cp LICENSE dist/LICENSE

test:
	go test -v ./...

lint:
	golangci-lint run

docker:
	docker build -t ${IMAGE_NAME}:${VERSION} -t ${IMAGE_NAME}:latest .