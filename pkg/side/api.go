package side

import (
	"context"
	"github.com/git-roll/monkey2/pkg/conf"
	"github.com/git-roll/monkey2/pkg/notify"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Car interface {
	Start(io.Writer)
	Kill()
	Done() <-chan error
}

type chanWriter []chan string

func (w chanWriter) Write(p []byte) (n int, err error) {
	for _, c := range w {
		c <- string(p)
	}

	n = len(p)
	return
}

func (w chanWriter) Close() {
	for _, c := range w {
		close(c)
	}
}

func NewCar() Car {
	if len(os.Args) < 3 {
		return newPlaceholder()
	}

	r := &Runner{
		done: make(chan error, 1),
	}
	var args []string
	if len(os.Args) > 3 {
		args = os.Args[3:]
	}

	r.proc = exec.CommandContext(context.Background(), os.Args[2], args...)
	r.proc.Env = os.Environ()
	pwd := conf.SideCarPWD()
	if len(pwd) > 0 {
		r.proc.Dir = pwd
	}

	setTermSig(r.proc)
	return r
}

type Runner struct {
	//stdout chanWriter
	//stderr chanWriter
	proc *exec.Cmd
	done chan error
}

func (r *Runner) Start(console io.Writer) {
	r.proc.Stdout = console
	r.proc.Stderr = console

	go func() {
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

	//r.stdout.Close()
	//r.stderr.Close()
}

//func (r *Runner) Stdout() <-chan string {
//	ch := make(chan string)
//	r.stdout = append(r.stdout, ch)
//	return ch
//}
//
//func (r *Runner) Stderr() <-chan string {
//	ch := make(chan string)
//	r.stderr = append(r.stderr, ch)
//	return ch
//}
