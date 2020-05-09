package char

import (
	"fmt"
	"time"
)

func Insane() Monkey {
	return &insaneMonkey{}
}

type insaneMonkey struct {
	idle *time.Timer
}

func (m insaneMonkey) Halt() {
	panic("implement me")
}

func (m *insaneMonkey) StartWork() {
	idle := randomIdleTime()
	fmt.Printf("☕️ break for %s", idle)
	m.idle = time.NewTimer(idle)
}
