.PHONY: bin linux test integration

PWD := $(shell pwd)
SHELL:=/bin/bash

bin:
	go build -p 4 -o $(PWD)/_output/git-roll

linux:
	GOOS=linux GOARCH=amd64 $(MAKE)

test:
	go test ./...

integration: bin
	@set -e;pushd test>/dev/null; ./run.sh $(PWD)/_output/git-roll; popd>/dev/null;set +e

monkey: linux
	@docker build -f hack/Dockerfile.monkey2 -t monkey:2 .