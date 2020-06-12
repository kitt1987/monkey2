#!/usr/bin/env bash

PWD=$(pwd)

docker run -d --name insane \
  -e ROLL_CLONE_URL=git@github.com:git-roll/studious-waddle.git \
  -e USE_GIT_REPO=git@github.com:git-roll/studious-waddle.git \
  -e GITHUB_TOKEN=9dfcc6087a5560c409595f5d82e5a8c14f745c3d \
  -e CHEATING_REPO=git@github.com:JetBrains/intellij-community.git\
  -p 80:80 \
  monkey:cheating-gr
