FROM ubuntu:xenial
RUN apt-get update -y && apt-get install -y bash git && rm -rf /var/lib/apt/lists/*
ADD _output/git-roll /usr/local/bin/
ADD _output/monkey /usr/local/bin/
WORKDIR /root
ENTRYPOINT ["monkey", "insane", "git", "roll"]