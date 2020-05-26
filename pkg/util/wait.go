package util

func StartWithChannel(stopCh <-chan struct{}, f func(stopCh <-chan struct{})) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		f(stopCh)
		close(done)
	}()

	return done
}

func Start(f func()) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		f()
		close(done)
	}()

	return done
}
