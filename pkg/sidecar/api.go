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

type Runner struct {
	stdout chanWriter
	stderr chanWriter
	ctx context.Context
	Kill context.CancelFunc
}

type RunnerStdout struct {

}

func New(name string, arg ...string) (r *Runner) {
	r = &Runner{}

	r.ctx, r.Kill = context.WithCancel(context.Background())

	gitProcess := exec.CommandContext(r.ctx, name, arg...)
	gitProcess.Stdout = r.stdout
	gitProcess.Stderr = r.stderr
	return
}