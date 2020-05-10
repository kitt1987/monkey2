package main

import (
	"fmt"
	"github.com/git-roll/monkey2/pkg/char"
	"github.com/git-roll/monkey2/pkg/conf"
	"k8s.io/apimachinery/pkg/util/wait"
	"os"
	"os/signal"
	"syscall"
)

const Usage = `monkey [name], the name could be [insane].
You can also run a sidecar to watch the monkey.`

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

	stopC := make(chan struct{})
	wg := wait.Group{}
	wg.StartWithChannel(stopC, monkey.StartWork)

	for {
		select {
		case <-signCh:
			signal.Stop(signCh)
			close(stopC)
		}
	}

	wg.Wait()
}
