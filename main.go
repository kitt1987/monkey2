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

	stopC := make(chan struct{})
	wg := wait.Group{}

	if wss != nil {
		wg.StartWithChannel(stopC, wss.Run)
	}

	wt := conf.Worktree()
	repo := conf.UseGitRepo()
	var panicRecovery func(string)
	if len(repo) > 0 {
		sideDupNotifier := newFilterNotifier(sideNotifier)
		monDupNotifier := newFilterNotifier(monNotifier)
		sideNotifier = sideDupNotifier
		monNotifier = monDupNotifier

		notify.Set(monNotifier)

		bootAt := time.Now()
		panicRecovery = func(msg string) {
			writeLastWordsToRepo(repo, wt, msg, monDupNotifier.LastNotes(), sideDupNotifier.LastNotes(), bootAt)
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

	sidecar := side.NewCar()
	sidecar.Start(sideNotifier)

	var monkey char.Monkey
	switch os.Args[1] {
	case "insane":
		notify.Printf("🐲 I'm a monkey. I'm INSANE!\n")
		monkey = char.Insane(wt, panicRecovery)
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
	fmt.Printf("RECORD LAST WORDS")
	y, m, d := boot.Date()
	h, min, s := boot.Clock()
	ts := fmt.Sprintf("%d%02d%02d-%02d%02d%02d", y, m, d, h, min, s)

	err := callGit(worktree, "checkout", "-B", "lastword"+ts, "master")
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
		lastwordsFile, []byte(strings.Join([]string{message, monkeyLog, sidecarLog}, "===\n")), 0644,
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

	err = callGit(worktree, "push", "origin", "master")
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

func newFilterNotifier(notifier io.WriteCloser) *noteFilter  {
	return &noteFilter{
		notifier:     notifier,
		lastNotes:    make([]string, 0, 50),
		maxLastNotes: 50,
	}
}

type noteFilter struct {
	notifier     io.WriteCloser
	lastNotes    []string
	maxLastNotes int
}

func (n2 *noteFilter) Write(p []byte) (n int, err error) {
	if len(n2.lastNotes) < n2.maxLastNotes {
		n2.lastNotes = append(n2.lastNotes, string(p))
	} else {
		n2.lastNotes[0] = ""
		n2.lastNotes = append(n2.lastNotes[1:], string(p))
	}

	return n2.notifier.Write(p)
}

func (n2 *noteFilter) Close() error {
	return n2.notifier.Close()
}

func (n2 noteFilter) LastNotes() string {
	return strings.Join(n2.lastNotes, "\n")
}
