package char

func Cheating(worktree, targetRepo string, recover func(string)) Monkey {
    return &cheatingMonkey{}
}

type cheatingMonkey struct {

}

func (c cheatingMonkey) StartWork(stopC <-chan struct{}) {
    panic("implement me")
}

func (c cheatingMonkey) Halt() {
    panic("implement me")
}
