VERSION := $(shell git describe --tags --dirty --always)
COMPILED := $(shell date +%s)

LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.compiled=$(COMPILED) -s -w"
IMAGE_NAME := davidsbond/mona
GOOS := ${shell go env GOOS}
GOARCH := ${shell go env GOARCH}

build:
	go build $(LDFLAGS) -o dist/mona
	cp README.md dist/README.md
	cp LICENSE dist/LICENSE

test:
	go test -v ./...

lint:
	golangci-lint run

docker-build: build
	docker build -t ${IMAGE_NAME}:${VERSION} -t ${IMAGE_NAME}:latest .

docker-push:
	docker push ${IMAGE_NAME}:${VERSION}
	docker push ${IMAGE_NAME}:latest