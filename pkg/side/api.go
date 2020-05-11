package side

import (
	"context"
	"fmt"
	"github.com/git-roll/monkey2/pkg/conf"
	"os"
	"os/exec"
	"strings"
)

type Car interface {
	Start()
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
	r.proc.Dir = conf.Worktree()
	return r
}

type Runner struct {
	//stdout chanWriter
	//stderr chanWriter
	proc *exec.Cmd
	std  *os.File
	done chan error
}

func (r *Runner) Start() {
	stdfile := conf.SidecarStdFile()
	f, err := os.OpenFile(stdfile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(fmt.Sprintf("%s:%s", stdfile, err))
	}

	r.std = f

	r.proc.Stdout = f
	r.proc.Stderr = f

	go func() {
		fmt.Printf(`🚁 Start sidecar "%s"`+"\n", strings.Join(r.proc.Args, " "))
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
	defer func() {
		if r.std != nil {
			r.std.Close()
		}
	}()

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
