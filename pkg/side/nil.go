package side

import "io"

type placeholder struct {
	done chan error
}

func (p *placeholder) Kill() {
	close(p.done)
}

func (p *placeholder) Start(io.Writer) {
}

func (p *placeholder) Done() <-chan error {
	return p.done
}

func newPlaceholder() *placeholder {
	return &placeholder{done: make(chan error)}
}
