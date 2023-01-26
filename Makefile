ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

all:
	go build -o build/ ${ROOT_DIR}/

test:
	go test -v ./...