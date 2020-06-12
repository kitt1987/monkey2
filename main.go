package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/git-roll/monkey2/pkg/char"
	"github.com/git-roll/monkey2/pkg/conf"
	"github.com/git-roll/monkey2/pkg/notify"
	"github.com/git-roll/monkey2/pkg/side"
	"github.com/git-roll/monkey2/pkg/ws"
	"k8s.io/apimachinery/pkg/util/wait"
)

const Usage = `monkey [name] [sidecar]

name could be one of [insane, cheating].
You can also run a sidecar to watch the monkey. e.g.

> monkey insane git roll

> monkey cheating git@github.com:git-roll/monkey2.git git roll`

func main() {
	if len(os.Args) < 2 {
		fmt.Println(Usage)
		return
	}

	signCh := make(chan os.Signal, 3)
	signal.Ignore(syscall.SIGPIPE)
	signal.Notify(signCh, syscall.SIGSEGV, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)

	var wss *ws.Server
	var sideNotifier, monNotifier io.WriteCloser
	if conf.WebSocketPort() > 0 {
		wss = ws.NewServer()
		sideNotifier = wss.SidecarNotifier()
		monNotifier = wss.MonkeyNotifier()
	} else {
		sideNotifier, monNotifier = createNotifier()
	}

	defer sideNotifier.Close()
	defer monNotifier.Close()

	stopC := make(chan struct{})
	wg := wait.Group{}

	if wss != nil {
		wg.StartWithChannel(stopC, wss.Run)
	}

	wt := conf.Worktree()
	repo := conf.UseGitRepo()
	var panicRecovery func(string)
	if len(repo) > 0 {
		sideDupNotifier := notify.NewMirrorNotifier(sideNotifier)
		monDupNotifier := notify.NewMirrorNotifier(monNotifier)
		sideNotifier = sideDupNotifier
		monNotifier = monDupNotifier

		notify.Set(monNotifier)

		bootAt := time.Now()
		panicRecovery = func(msg string) {
			writeLastWordsToRepo(repo, wt, msg, monDupNotifier.JoinNotices(), sideDupNotifier.JoinNotices(), bootAt)
		}

		notify.Printf("🚁 Clone %s=>%s\n", repo, wt)
		out, err := exec.Command("git", "clone", repo, wt).Output()
		if err != nil {
			fmt.Printf("Can't clone from %s: %s\n%s", repo, err.Error(), string(out))
			return
		}
	} else {
		notify.Set(monNotifier)
	}

	var sidecarArgs []string
	var monkey char.Monkey
	switch os.Args[1] {
	case "insane":
		notify.Printf("🐲 I'm a monkey. I'm INSANE!\n")
		monkey = char.Insane(wt, panicRecovery)
		sidecarArgs = os.Args[2:]
	case "cheating":
		if len(os.Args) < 3 {
			fmt.Println(Usage)
			return
		}

		notify.Printf("🦊 I'm a monkey. I'm going to cheat some repos!\n")
		monkey = char.Cheating(wt, os.Args[2], panicRecovery)
		sidecarArgs = os.Args[3:]
	default:
		fmt.Println(Usage)
		return
	}

	sidecar := side.NewCar(sidecarArgs, panicRecovery)
	sidecar.Start(sideNotifier)

	wg.StartWithChannel(stopC, monkey.StartWork)

	for {
		select {
		case <-sidecar.Done():
			notify.Printf("🩸 Sidecar broke!\n")
			signal.Stop(signCh)
			notify.Printf("🛎 Monkey exit!\n")
			close(stopC)
			wg.Wait()
			return

		case sig, ok := <-signCh:
			if !ok {
				return
			}

			notify.Printf("🛎 Got signal %s\n", sig)
			signal.Stop(signCh)
			notify.Printf("🛎 Monkey exit!\n")
			close(stopC)
			wg.Wait()
			notify.Printf("🛎 Stop sidecar!\n")
			sidecar.Kill()
			return
		}
	}
}

func createNotifier() (side, monkey io.WriteCloser) {
	stdfile := conf.SidecarStdFile()
	if len(stdfile) > 0 {
		f, err := os.OpenFile(stdfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			panic(fmt.Sprintf("%s:%s", stdfile, err))
		}

		return f, os.Stdout
	}

	return os.Stdout, os.Stdout
}

func writeLastWordsToRepo(repo, worktree, message, monkeyLog, sidecarLog string, boot time.Time) {
	fmt.Println("RECORD LAST WORDS")
	y, m, d := boot.Date()
	h, min, s := boot.Clock()
	ts := fmt.Sprintf("%d%02d%02d-%02d%02d%02d", y, m, d, h, min, s)

	branch := "lastword"+ts
	err := callGit(worktree, "checkout", "-B", branch, "master")
	if err != nil {
		fmt.Printf("checkout: %s", err)
		return
	}

	lastwordsDir := filepath.Join(worktree, ".lastwords")
	err = os.MkdirAll(lastwordsDir, 0755)
	if err != nil {
		fmt.Printf("mkdir: %s", err)
		return
	}

	lastwordsFile := filepath.Join(lastwordsDir, ts)
	err = ioutil.WriteFile(
		lastwordsFile,
		[]byte(strings.Join([]string{message, monkeyLog, sidecarLog}, "🤖🤖🤖🤖🤖🤖🤖🤖🤖🤖🤖🤖🤖🤖🤖🤖🤖🤖\n")),
		0644,
	)
	if err != nil {
		fmt.Printf("write: %s", err)
		return
	}

	err = callGit(worktree, "add", lastwordsFile)
	if err != nil {
		fmt.Printf("add: %s", err)
		return
	}

	err = callGit(worktree, "commit", "-m", "last words")
	if err != nil {
		fmt.Printf("commit: %s", err)
		return
	}

	err = callGit(worktree, "push", "origin", branch)
	if err != nil {
		fmt.Printf("push: %s", err)
		return
	}
}

func callGit(worktree string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = worktree
	return cmd.Run()
}
