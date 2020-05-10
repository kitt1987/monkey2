package main

import (
	"fmt"
	"github.com/git-roll/monkey2/pkg/char"
	"github.com/git-roll/monkey2/pkg/conf"
	"github.com/git-roll/monkey2/pkg/sidecar"
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
		monkey = char.Insane(conf.Worktree())
	default:
		fmt.Println(Usage)
		return
	}

	var sc *sidecar.Runner
	if len(os.Args) > 2 {
		sc = sidecar.New(os.Args[2], os.Args[2:]...)
		err := sc.Start()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	stopC := make(chan struct{})
	wg := wait.Group{}
	wg.StartWithChannel(stopC, monkey.StartWork)
	defer func() {
		wg.Wait()
		if sc != nil {
			sc.Kill()
		}
	}()

	for {
		select {
		case <-signCh:
			signal.Stop(signCh)
			close(stopC)
			return
		}
	}
}
