.PHONY: bin

PWD := $(shell pwd)
SHELL:=/bin/bash

bin:
	go build -p 4 -o $(PWD)/_output/monkey

linux:
	GOOS=linux GOARCH=amd64 $(MAKE)

insane: linux
	@docker build -f hack/Dockerfile.insane -t monkey:insane .
