#!/usr/bin/env bash

PWD=$(pwd)

docker run -d --name insane \
  -e ROLL_CLONE_URL=git@github.com:git-roll/expert-garbanzo.git \
  -e GITHUB_TOKEN=50c2351a3de25fc1bf1147bebacd91a56f48cea8 \
  -p 80:80 \
  monkey:insane-gr
