FROM monkey:insane

RUN apt-get update -y && apt-get install -y hub && rm -rf /var/lib/apt/lists/*
ENV ROLL_GIT_USER_NAME=insane-monkey
ENV ROLL_GIT_USER_EMAIL=insane.monkey@releases.fyi
ENV EXCLUDED_FILES=README.md
ENV WEBSOCKET_PORT="80"
