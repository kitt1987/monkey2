package conf

import (
	"github.com/git-roll/monkey2/pkg/notify"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	EnvCoffeeTimeUpperBound = "COFFEE_TIME"
	EnvNameLength           = "NAME_LENGTH"
	EnvWriteOnceLength      = "LENGTH_WRITE_ONCE"
	EnvPercentageFileOP     = "PERCENTAGE_FILE_OPERATION"
	EnvPercentageCmd        = "PERCENTAGE_CMD"
	EnvWorktree             = "WORKTREE"
	EnvSidecarStdFile       = "SIDECAR_STD_FILE"
	EnvSidecarPWD           = "SIDECAR_PWD"
	EnvWebSocketPort        = "WEBSOCKET_PORT"
	EnvExcludedFiles        = "EXCLUDED_FILES"
	EnvCmdSeqFile           = "CMD_SEQ_FILE"
	EnvCmdBuildPR           = "CMD_BUILD_PR"
	EnvGitRepo              = "USE_GIT_REPO"
	EnvCheatingRepo         = "CHEATING_REPO"
)

var noticeOnce = make(map[string]bool)

func notice(key string, hint string, v interface{}) {
	if !noticeOnce[key] {
		notify.Printf(hint+`. Set environment variable "%s" to change.`+"\n", v, key)
		noticeOnce[key] = true
	}
}

func CoffeeTimeUpperBound() string {
	v := os.Getenv(EnvCoffeeTimeUpperBound)
	if len(v) == 0 {
		v = (15 * time.Second).String()
	}

	notice(EnvCoffeeTimeUpperBound, `🚁 Coffee time would be up to %s`, v)
	return v
}

func Worktree() (wt string) {
	wt = os.Getenv(EnvWorktree)
	if len(wt) == 0 {
		wt = filepath.Join(os.TempDir(), "monkey")
	}

	notice(EnvWorktree, `🚁 The workdir will be "%s"`, wt)
	return
}

func UseGitRepo() (repo string) {
	repo = os.Getenv(EnvGitRepo)
	if len(repo) == 0 {
		return
	}

	notice(EnvGitRepo, `🚁 Use Git repo "%s"`, repo)
	return
}

func SideCarPWD() string {
	wt := os.Getenv(EnvSidecarPWD)
	if len(wt) == 0 {
		return ""
	}

	notice(EnvSidecarPWD, `🚁 The workdir of SideCar will be "%s"`, wt)
	return wt
}

func SidecarStdFile() string {
	std := os.Getenv(EnvSidecarStdFile)
	if len(std) == 0 {
		std = filepath.Join(os.TempDir(), "monkey.sidecar")
	}

	notice(EnvSidecarStdFile, `🚁 Stdout of the sidecar will be written to "%s"`, std)
	return std
}

func CmdSeqFile() string {
	seq := os.Getenv(EnvCmdSeqFile)
	if len(seq) > 0 {
		notice(EnvCmdSeqFile, `🚁 Use command sequence in file "%s"`, seq)
	}

	return seq
}

func CmdBuildPR() string {
	pr := os.Getenv(EnvCmdBuildPR)
	if len(pr) > 0 {
		notice(EnvCmdBuildPR, `🚁 Use command "%s" to build PR`, pr)
	}

	return pr
}

func CheatingRepo() string {
	repo := os.Getenv(EnvCheatingRepo)
	if len(repo) > 0 {
		notice(EnvCheatingRepo, `🚁 Cheat repo "%s"`, repo)
	}

	return repo
}

func WebSocketPort() int {
	return envInt(EnvWebSocketPort, 80,
		`🚁 Logs are exposed on port %d`,
	)
}

func NameLength() int {
	return envInt(EnvNameLength, 8, `🚁 Length of file/dir name would be %d`)
}

func WriteOnceLengthUpperBound() int {
	return envInt(
		EnvWriteOnceLength, 512,
		`🚁 Length of each file write would be %d`,
	)
}

func PercentageFileOP() int {
	return envInt(EnvPercentageFileOP, 70,
		`🚁 %d%% filesystem operations would be on files`,
	)
}

func PercentageCmd() int {
	return envInt(EnvPercentageCmd, 10,
		`🚁 %d%% command would be executed against file operations`,
	)
}

func ExcludedFiles() map[string]bool {
	v := os.Getenv(EnvExcludedFiles)
	if len(v) == 0 {
		return nil
	}

	ex := make(map[string]bool)
	for _, f := range strings.Split(v, ",") {
		ex[f] = true
	}

	return ex
}

func envInt(key string, def int, hint string) (i int) {
	v := os.Getenv(key)
	if len(v) > 0 {
		i64, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			panic(os.Getenv(key))
		}
		i = int(i64)
	} else {
		i = def
	}

	notice(key, hint, i)
	return
}
