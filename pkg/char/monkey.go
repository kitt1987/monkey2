package char

type Monkey interface {
	StartWork(stopC <-chan struct{})
	Halt()
}
