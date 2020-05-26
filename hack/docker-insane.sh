#!/usr/bin/env bash

PWD=$(pwd)

docker run -d --name insane \
  -e ROLL_CLONE_URL=git@github.com:git-roll/expert-garbanzo.git \
  -e ROLL_GIT_USER_NAME=insane-monkey \
  -e ROLL_GIT_USER_EMAIL=insane.monkey@releases.fyi \
  -e EXCLUDED_FILES=README.md \
  -e WEBSOCKET_PORT="80" \
  -v "${PWD}/key/insane-monkey":/root/.ssh/id_rsa \
  -v "${PWD}/key/insane-monkey.pub":/root/.ssh/id_rsa.pub \
  -v "${PWD}/key/ssh-config":/root/.ssh/config \
  -p 80:80 \
  monkey:insane
