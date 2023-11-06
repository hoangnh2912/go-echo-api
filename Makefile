#@IgnoreInspection BashAddShebang
export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export DEBUG=true
export APP=go-echo-api
export LDFLAGS="-w -s"

all: build test

build:
	go build -race  .

build-static:
	CGO_ENABLED=0 go build -race -v -o $(APP) -a -installsuffix cgo -ldflags $(LDFLAGS) .

run:
	go run -race .

dev:
	bash -c "trap 'docker-compose down' EXIT; docker-compose -f ./docker-compose.dev.yml up -d --build && air"

swagger:
	swag i
############################################################
# Test
############################################################

test:
	go test -v -race ./...

.PHONY: build run build-static test container
