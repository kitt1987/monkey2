package sidecar

import (
	"context"
	"os/exec"
)

type chanWriter []chan string

func (w chanWriter) Write(p []byte) (n int, err error) {
	for _, c := range w{
		c <- string(p)
	}

	n = len(p)
	return
}

func New(name string, arg ...string) (r *Runner) {
	r = &Runner{}

	r.ctx, r.Kill = context.WithCancel(context.Background())
	r.proc = exec.CommandContext(r.ctx, name, arg...)
	return
}

type Runner struct {
	stdout chanWriter
	stderr chanWriter
	ctx context.Context
	Kill context.CancelFunc
	proc *exec.Cmd
}

func (r *Runner) Start() error {
	r.proc.Stdout = r.stdout
	r.proc.Stderr = r.stderr
	return r.proc.Start()
}

func (r *Runner) Stdout() <-chan string {
	ch := make(chan string)
	r.stdout = append(r.stdout, ch)
	return ch
}

func (r *Runner) Stderr() <-chan string {
	ch := make(chan string)
	r.stderr = append(r.stderr, ch)
	return ch
}
