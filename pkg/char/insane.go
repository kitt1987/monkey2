package char

import (
	"fmt"
	"github.com/git-roll/monkey2/pkg/conf"
	"github.com/git-roll/monkey2/pkg/op"
	"time"
)

func Insane(worktree string) Monkey {
	m := &insaneMonkey{
		worktree: op.NewWorktree(worktree),
	}
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
	obBias := NewObjectBias()
	obBias.Set(int(op.FSFile), conf.PercentageFileOP())
	obBias.Set(int(op.FSDir), 100 - conf.PercentageFileOP())

	allDirs := m.worktree.AllDirs()
	dirOpBias := NewDirOPBias()
	if len(allDirs) == 0 {
		dirOpBias.Set(int(op.FSCreate), 100)
	} else {
		dirOpBias.Set(int(op.FSCreate), 34)
		dirOpBias.Set(int(op.FSDelete), 33)
		dirOpBias.Set(int(op.FSRename), 33)
	}

	allFiles := m.worktree.AllFiles()
	fileOpBias := NewFileOPBias()
	if len(allFiles) == 0 {
		fileOpBias.Set(int(op.FSCreate), 100)
	} else {
		fileOpBias.Set(int(op.FSCreate), 25)
		fileOpBias.Set(int(op.FSDelete), 25)
		fileOpBias.Set(int(op.FSRename), 25)
		fileOpBias.Set(int(op.FSOverride), 25)
	}

	ob, op := randomFSOp(obBias, dirOpBias, fileOpBias)
	m.worktree.Apply(ob, op, m.prepareArgs(allFiles, allDirs))
}

func (m *insaneMonkey) prepareArgs(allFiles, allDirs []string) *op.WorktreeOPArgs {
	f := randomItem(allFiles)
	size := m.worktree.FileSize(f)
	offset := randomN64(size)

	return &op.WorktreeOPArgs{
		NewRelativeFilePath:     "f-" + randomText(conf.NameLength()),
		ExistedRelativeFilePath: f,
		NewRelativeDirPath:      "d-" + randomText(conf.NameLength()),
		ExistedRelativeDirPath:  randomItem(allDirs),
		Content:                 randomText(randomN(conf.WriteOnceLengthUpperBound())),
		Offset:                  offset,
		Size:                    randomN64(size - offset),
	}
}
