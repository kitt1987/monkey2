FROM monkey:insane

ADD _output/git-roll /usr/local/bin/
ADD hack/insane-gr/hub/hub /usr/local/bin
ENV ROLL_GIT_USER_NAME=insane-monkey
ENV ROLL_GIT_USER_EMAIL=insane.monkey@releases.fyi
ENV WEBSOCKET_PORT="80"
ENV EXCLUDED_FILES=README.md,git-roll.yml

WORKDIR /root
ADD key/insane-monkey .ssh/id_rsa
ADD key/insane-monkey.pub .ssh/id_rsa.pub
ADD key/ssh-config .ssh/config
ADD hack/insane-gr/monkey-cmd .
ENV CMD_SEQ_FILE=/root/monkey-cmd
ENV WORKTREE="/root/monkey_work"
ENV ROLL_WORKTREE=${WORKTREE}

ENTRYPOINT ["monkey", "insane", "git", "roll"]
