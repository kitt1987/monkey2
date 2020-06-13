package char

import (
	"github.com/git-roll/monkey2/pkg/notify"
	"os"
	"time"
)

type Monkey interface {
	StartWork(stopC <-chan struct{})
	Halt()
}

type monkeyChar interface {
	Work()
}

type monkey struct {
	recover  func(string)
	idle     *time.Timer
	monkeyChar monkeyChar
}

func (m monkey) Halt() {
	if m.idle != nil && !m.idle.Stop() {
		<-m.idle.C
	}
}

func (m *monkey) StartWork(stopC <-chan struct{}) {
	if m.recover != nil {
		defer func() {
			if r := recover(); r != nil {
				if msg, ok := r.(string); ok {
					m.recover(msg)
				}

				os.Exit(2)
			}
		}()
	}

	m.idle = time.NewTimer(0)
	for {
		select {
		case <-m.idle.C:
			m.monkeyChar.Work()

			idle := randomCoffeeTime()
			notify.Printf("☕️ coffee time: %s\n", idle)
			m.idle.Reset(idle)
		case <-stopC:
			m.Halt()
			return
		}
	}
}
