FROM monkey

ADD _output/git-roll /usr/local/bin/
ADD hack/insane-gr/hub/hub /usr/local/bin
ENV ROLL_GIT_USER_NAME=cheating-monkey
ENV ROLL_GIT_USER_EMAIL=cheating.monkey@releases.fyi
ENV WEBSOCKET_PORT="80"
ENV EXCLUDED_FILES=README.md,git-roll.yml

WORKDIR /root
ADD key/cheating-monkey .ssh/id_rsa
ADD key/cheating-monkey.pub .ssh/id_rsa.pub
ADD key/ssh-config .ssh/config
ADD hack/insane-gr/monkey-cmd .
ENV CMD_SEQ_FILE=/root/monkey-cmd
ENV WORKTREE="/root/monkey_work"
ENV ROLL_WORKTREE=${WORKTREE}
ENV CMD_BUILD_PR="git roll complete"

ENTRYPOINT ["monkey", "cheating", "git", "roll"]
