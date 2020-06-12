.PHONY: bin

PWD := $(shell pwd)
SHELL:=/bin/bash

bin:
	go build -p 4 -o $(PWD)/_output/monkey

linux:
	GOOS=linux GOARCH=amd64 $(MAKE)

monkey: linux
	@docker build -f hack/Dockerfile -t monkey:latest .

insane-gr:
	@docker build -f hack/gr/insane.dockerfile -t monkey:insane-gr .

cheating-gr:
	@docker build -f hack/gr/cheating.dockerfile -t monkey:cheating-gr .
