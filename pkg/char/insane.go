package char

import (
	"fmt"
	"github.com/git-roll/monkey2/pkg/op"
	"time"
)

func Insane() Monkey {
	m := &insaneMonkey{}
	return m
}

type insaneMonkey struct {
	idle     *time.Timer
	worktree *op.Worktree
}

func (m insaneMonkey) Halt() {
	if m.idle != nil && !m.idle.Stop() {
		<-m.idle.C
	}
}

func (m *insaneMonkey) StartWork(stopC <-chan struct{}) {
	m.idle = time.NewTimer(time.Nanosecond)
	for {
		select {
		case _, ok := <-m.idle.C:
			if !ok {
				return
			}

			m.work()

			idle := randomCoffeeTime()
			fmt.Printf("☕️ coffee time: %s", idle)
			m.idle.Reset(idle)
		case <-stopC:
			m.Halt()
		}
	}
}

func (m *insaneMonkey) work() {

}
