package main

import (
	"fmt"
	"github.com/git-roll/monkey2/pkg/char"
	"github.com/git-roll/monkey2/pkg/conf"
	"github.com/git-roll/monkey2/pkg/side"
	"k8s.io/apimachinery/pkg/util/wait"
	"os"
	"os/signal"
	"syscall"
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
	signal.Notify(signCh, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)

	var monkey char.Monkey

	switch os.Args[1] {
	case "insane":
		fmt.Printf("🐲 I'm a monkey. I'm INSANE!\n")
		monkey = char.Insane(conf.Worktree())
	default:
		fmt.Println(Usage)
		return
	}

	sidecar := side.NewCar()
	sidecar.Start()

	stopC := make(chan struct{})
	wg := wait.Group{}
	wg.StartWithChannel(stopC, monkey.StartWork)

	for {
		select {
		case err, _ := <-sidecar.Done():
			if err != nil {
				fmt.Printf("🩸 Sidecar broke!\n")
			}

			signal.Stop(signCh)
			close(stopC)
			wg.Wait()
			return

		case _, ok := <-signCh:
			if !ok {
				return
			}

			signal.Stop(signCh)
			close(stopC)
			wg.Wait()
			sidecar.Kill()
			return
		}
	}
}
