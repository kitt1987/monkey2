package main

import (
	"fmt"
	"github.com/git-roll/monkey2/pkg/char"
	"github.com/git-roll/monkey2/pkg/conf"
	"github.com/git-roll/monkey2/pkg/notify"
	"github.com/git-roll/monkey2/pkg/side"
	"github.com/git-roll/monkey2/pkg/ws"
	"io"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/wait"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

const Usage = `monkey [name] [sidecar]

name could be one of [insane].
You can also run a sidecar to watch the monkey. e.g.

> monkey insane git roll`

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

	notify.Set(monNotifier)

	stopC := make(chan struct{})
	wg := wait.Group{}

	if wss != nil {
		wg.StartWithChannel(stopC, wss.Run)
	}

	wt := conf.Worktree()
	repo := conf.UseGitRepo()
	if len(repo) > 0 {
		notify.Printf("🚁 Clone %s=>%s\n", repo, wt)
		out, err := exec.Command("git", "clone", repo, wt).Output()
		if err != nil {
			fmt.Printf("Can't clone from %s: %s\n%s", repo, err.Error(), string(out))
			return
		}

		bootAt := time.Now()
		defer func() {
			if r := recover(); r != nil {
				if msg, ok := r.(string); ok {
					writeLastWordsToRepo(repo, wt, msg, monkey, sidecar, bootAt)
				}

				os.Exit(128)
			}
		}()
	}

	sidecar := side.NewCar()
	sidecar.Start(sideNotifier)

	var monkey char.Monkey
	switch os.Args[1] {
	case "insane":
		notify.Printf("🐲 I'm a monkey. I'm INSANE!\n")
		monkey = char.Insane(wt)
	default:
		fmt.Println(Usage)
		return
	}

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

			notify.Printf("🛎 Got signal %s", sig)
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
	y, m, d := boot.Date()
	h, min, s := boot.Clock()
	ts := fmt.Sprintf("%d%02d%02d-%02d%02d%02d", y, m, d, h, min, s)

	err := callGit(worktree, "checkout", "-B", "lastword"+ts, "master")
	if err != nil {
		return
	}

	lastwordsDir := filepath.Join(repo, ".lastwords")
	err = os.MkdirAll(lastwordsDir, 0755)
	if err != nil {
		return
	}

	lastwordsFile := filepath.Join(lastwordsDir, ts)
	err = ioutil.WriteFile(
		lastwordsFile, []byte(strings.Join([]string{message, monkeyLog, sidecarLog}, "---/n")), 0644,
		)
	if err != nil {
		return
	}

	err = callGit(worktree, "add", lastwordsFile)
	if err != nil {
		return
	}

	err = callGit(worktree, "commit", "-m", "last words")
	if err != nil {
		return
	}

	err = callGit(worktree, "push", "origin", "master")
	if err != nil {
		return
	}
}

func callGit(worktree string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = worktree
	return cmd.Run()
}

type noteFilter struct {
	notifier io.WriteCloser
	lastNotes []string
}

func (n2 noteFilter) Write(p []byte) (n int, err error) {
	return n2.notifier.Write(p)
}

func (n2 noteFilter) Close() error {
	return n2.notifier.Close()
}
