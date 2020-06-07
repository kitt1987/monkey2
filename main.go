package main

import (
	"fmt"
	"github.com/git-roll/monkey2/pkg/char"
	"github.com/git-roll/monkey2/pkg/conf"
	"github.com/git-roll/monkey2/pkg/notify"
	"github.com/git-roll/monkey2/pkg/side"
	"github.com/git-roll/monkey2/pkg/ws"
	"io"
	"k8s.io/apimachinery/pkg/util/wait"
	"os"
	"os/signal"
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

	sidecar := side.NewCar()
	sidecar.Start(sideNotifier)

	stopC := make(chan struct{})
	wg := wait.Group{}

	if wss != nil {
		wg.StartWithChannel(stopC, wss.Run)
	}

	wt, generated := conf.Worktree()
	if !generated {
		notify.Printf("Wait for sidecar to create the worktree dir\n")
		for {
			_, err := os.Lstat(wt)
			if os.IsNotExist(err) {
				time.Sleep(500 * time.Millisecond)
				continue
			}

			if err != nil {
				fmt.Println("Can't use the specified worktree ", wt, err.Error())
				return
			}

			break
		}
	}

	repo := conf.UseGitRepo()
	if len(repo) > 0 {
		notify.Printf("Wait for sidecar to create the worktree dir\n")
	}

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
