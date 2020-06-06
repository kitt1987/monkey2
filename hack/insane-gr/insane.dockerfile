FROM monkey:insane

ADD _output/git-roll /usr/local/bin/
ADD hack/insane-gr/hub/bin/hub /usr/local/bin
ENV ROLL_GIT_USER_NAME=insane-monkey
ENV ROLL_GIT_USER_EMAIL=insane.monkey@releases.fyi
ENV WEBSOCKET_PORT="80"

WORKDIR /root
ADD key/insane-monkey .ssh/id_rsa
ADD key/insane-monkey.pub .ssh/id_rsa.pub
ADD key/ssh-config .ssh/config
ADD monkey-cmd .
ENV CMD_SEQ_FILE=/root/monkey-cmd
ENV WORKTREE="/root/monkey_work"
WORKDIR ${WORKTREE}
ADD git-roll.yml .
ENV EXCLUDED_FILES=README.md,git-roll.yml

ENTRYPOINT ["monkey", "insane", "git", "roll"]
