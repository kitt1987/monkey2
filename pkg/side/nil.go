package side

type placeholder struct {
	done chan error
}

func (p *placeholder) Kill() {
	close(p.done)
}

func (p *placeholder) Start() error {
	return nil
}

func (p *placeholder) Done() <-chan error {
	return p.done
}

func newPlaceholder() *placeholder {
	return &placeholder{done: make(chan error)}
}