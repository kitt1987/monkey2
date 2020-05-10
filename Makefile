.PHONY: bin

PWD := $(shell pwd)
SHELL:=/bin/bash

bin:
	go build -p 4 -o $(PWD)/_output/monkey

insane:
	@docker build -f hack/Dockerfile.insane -t monkey:insane .
