VERSION := $(shell git describe --tags --dirty --always)
LDFLAGS := -ldflags '-X "main.version=$(VERSION)" -s -w'
IMAGE_NAME := davidsbond/mona
GOOS := ${shell go env GOOS}
GOARCH := ${shell go env GOARCH}

build:
	go build $(LDFLAGS) -o dist/mona
	cp README.md dist/README.md
	cp LICENSE dist/LICENSE
	mkdir -p release/${VERSION}
	tar -zcvf release/${VERSION}/mona_${VERSION}_${GOOS}_${GOARCH}.tar.gz dist 

test:
	go test -v ./...

lint:
	golangci-lint run

docker: build
	docker build -t ${IMAGE_NAME}:${VERSION} -t ${IMAGE_NAME}:latest .

docker-push:
	docker push ${IMAGE_NAME}:${VERSION}
	docker push ${IMAGE_NAME}:latest