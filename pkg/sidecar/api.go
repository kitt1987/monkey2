package sidecar

import (
	"context"
	"fmt"
	"github.com/git-roll/monkey2/pkg/conf"
	"os"
	"os/exec"
)

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

func New(name string, arg ...string) (r *Runner) {
	r = &Runner{}
	r.ctx, r.kill = context.WithCancel(context.Background())
	r.proc = exec.CommandContext(r.ctx, name, arg...)
	r.proc.Env = os.Environ()
	r.proc.Dir = conf.Worktree()
	return
}

type Runner struct {
	//stdout chanWriter
	//stderr chanWriter
	ctx    context.Context
	kill   context.CancelFunc
	proc   *exec.Cmd
	std *os.File
}

func (r *Runner) Start() error {
	stdfile := conf.SidecarStdFile()
	f, err := os.OpenFile(stdfile, os.O_RDWR, 0666)
	if err != nil {
		panic(fmt.Sprintf("%s:%s", stdfile, err))
	}

	r.std = f

	r.proc.Stdout = f
	r.proc.Stderr = f
	return r.proc.Start()
}

func (r *Runner) Kill() {
	r.kill()
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
