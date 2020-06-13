package side

import (
	"context"
	"github.com/git-roll/monkey2/pkg/conf"
	"github.com/git-roll/monkey2/pkg/notify"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Car interface {
	Start(io.Writer)
	Kill()
	Done() <-chan error
}

func NewCar(args []string, recover func(string)) Car {
	if len(os.Args) < 3 {
		return newPlaceholder()
	}

	r := &Runner{
		recover: recover,
		done: make(chan error, 1),
	}

	r.proc = exec.CommandContext(context.Background(), args[0], args[1:]...)
	r.proc.Env = os.Environ()
	pwd := conf.SideCarPWD()
	if len(pwd) > 0 {
		r.proc.Dir = pwd
	}

	setTermSig(r.proc)
	return r
}

type Runner struct {
	recover func(string)
	proc *exec.Cmd
	done chan error
}

func (r *Runner) Start(console io.Writer) {
	r.proc.Stdout = console
	r.proc.Stderr = console

	go func() {
		defer func() {
			if re := recover(); re != nil {
				if msg, ok := re.(string); ok {
					stackBuf := make([]byte, 8192)
					numBytes := runtime.Stack(stackBuf, false)
					r.recover(msg + "\n" + string(stackBuf[:numBytes]))
				}

				os.Exit(2)
			}
		}()

		notify.Printf(`🚁 Start sidecar "%s"`+"\n", strings.Join(r.proc.Args, " "))
		r.done <- r.proc.Run()
		close(r.done)
	}()
}

func (r *Runner) Done() <-chan error {
	return r.done
}

func (r *Runner) Kill() {
	if r.proc.Process == nil {
		return
	}

	r.proc.Process.Kill()
	if err := r.proc.Wait(); err != nil {
		return
	}
}
