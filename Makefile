# Build information
VERSION := $(shell git describe --tags --dirty --always)
COMPILED := $(shell date +%s)
INSTALL_DIR := $(shell go env GOPATH)/bin

# Flags to pass to go build
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.compiled=$(COMPILED) -s -w"

# Docker image info
IMAGE_NAME := davidsbond/mona

# Builds the binary, injecting version & compile time information, creates a
# dist/ directory containing the binary, license and readme.
build:
	go build $(LDFLAGS) -o dist/mona
	cp README.md dist/README.md
	cp LICENSE dist/LICENSE

# Runs all tests
test:
	go test -v ./... -bench=. -race

# Lints all go packages
lint:
	golangci-lint run

# Executes the 'build' recipe, then generates a docker image for the latest and
# current version tags.
docker-build: build
	docker build -t ${IMAGE_NAME}:${VERSION} -t ${IMAGE_NAME}:latest .

# Pushes the current version of the docker image under the version and latest tags
docker-push:
	docker push ${IMAGE_NAME}:${VERSION}
	docker push ${IMAGE_NAME}:latest

install: build
	sudo cp dist/mona $(INSTALL_DIR)/mona