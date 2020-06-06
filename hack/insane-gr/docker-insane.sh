#!/usr/bin/env bash

PWD=$(pwd)

docker run -d --name insane \
  -e ROLL_CLONE_URL=git@github.com:git-roll/expert-garbanzo.git \
  -p 80:80 \
  monkey:insane-gr
