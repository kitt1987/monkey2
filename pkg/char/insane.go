package char

import (
	"github.com/git-roll/monkey2/pkg/cmd"
	"github.com/git-roll/monkey2/pkg/conf"
	"github.com/git-roll/monkey2/pkg/fs"
	"github.com/git-roll/monkey2/pkg/notify"
	"time"
)

func Insane(worktree string) Monkey {
	m := &insaneMonkey{
		worktree: fs.NewWorktree(worktree),
	}

	seq := conf.CmdSeqFile()
	if len(seq) > 0 {
		m.commands = cmd.NewSeq(seq)
	}

	return m
}

type insaneMonkey struct {
	idle     *time.Timer
	worktree fs.Worktree
	commands *cmd.Seq
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
		case <-m.idle.C:
			m.work()

			idle := randomCoffeeTime()
			notify.Printf("☕️ coffee time: %s\n", idle)
			m.idle.Reset(idle)
		case <-stopC:
			m.Halt()
			return
		}
	}
}

func (m *insaneMonkey) work() {
	if m.commands == nil {
		m.fsWork()
		return
	}

	bias := NewActivityBias()
	bias.Set(int(CMDActivity), conf.PercentageCmd())
	bias.Set(int(FSActivity), 100-conf.PercentageCmd())
	activity := MonkeyActivity(bias.RandomObject())
	switch activity {
	case CMDActivity:
		m.cmdWork()
	case FSActivity:
		m.fsWork()
	default:
		panic(activity)
	}
}

func (m *insaneMonkey) cmdWork() {
	obBias := NewObjectBias()
	obBias.Set(int(fs.File), conf.PercentageFileOP())
	obBias.Set(int(fs.Dir), 100-conf.PercentageFileOP())

	allDirs := m.worktree.AllDirs()
	dirOpBias := NewDirOPBias()
	if len(allDirs) == 0 {
		dirOpBias.Set(int(fs.Create), 100)
	} else {
		dirOpBias.Set(int(fs.Create), 34)
		dirOpBias.Set(int(fs.Delete), 33)
		dirOpBias.Set(int(fs.Rename), 33)
	}

	allFiles := m.worktree.AllFiles()
	fileOpBias := NewFileOPBias()
	if len(allFiles) == 0 {
		fileOpBias.Set(int(fs.Create), 100)
	} else {
		fileOpBias.Set(int(fs.Create), 20)
		fileOpBias.Set(int(fs.Delete), 20)
		fileOpBias.Set(int(fs.Rename), 20)
		fileOpBias.Set(int(fs.Override), 40)
	}

	ob, op := randomFSOp(obBias, fileOpBias, dirOpBias)
	m.worktree.Apply(ob, op, m.prepareArgs(allFiles, allDirs))
}

func (m *insaneMonkey) fsWork() {
	obBias := NewObjectBias()
	obBias.Set(int(fs.File), conf.PercentageFileOP())
	obBias.Set(int(fs.Dir), 100-conf.PercentageFileOP())

	allDirs := m.worktree.AllDirs()
	dirOpBias := NewDirOPBias()
	if len(allDirs) == 0 {
		dirOpBias.Set(int(fs.Create), 100)
	} else {
		dirOpBias.Set(int(fs.Create), 34)
		dirOpBias.Set(int(fs.Delete), 33)
		dirOpBias.Set(int(fs.Rename), 33)
	}

	allFiles := m.worktree.AllFiles()
	fileOpBias := NewFileOPBias()
	if len(allFiles) == 0 {
		fileOpBias.Set(int(fs.Create), 100)
	} else {
		fileOpBias.Set(int(fs.Create), 20)
		fileOpBias.Set(int(fs.Delete), 20)
		fileOpBias.Set(int(fs.Rename), 20)
		fileOpBias.Set(int(fs.Override), 40)
	}

	ob, op := randomFSOp(obBias, fileOpBias, dirOpBias)
	m.worktree.Apply(ob, op, m.prepareArgs(allFiles, allDirs))
}

func (m *insaneMonkey) prepareArgs(allFiles, allDirs []string) *fs.WorktreeOPArgs {
	args := &fs.WorktreeOPArgs{
		NewRelativeFilePath: "f-" + randomName(conf.NameLength()),
		NewRelativeDirPath:  "d-" + randomName(conf.NameLength()),
		Content:             randomText(randomN(conf.WriteOnceLengthUpperBound())),
	}

	if len(allFiles) > 0 {
		args.ExistedRelativeFilePath = randomItem(allFiles)
		size := m.worktree.FileSize(args.ExistedRelativeFilePath)
		if size == 0 {
			args.Offset = 0
			args.Size = 0
		} else {
			args.Offset = randomN64(size)
			args.Size = randomN64(size - args.Offset)
		}
	}

	if len(allDirs) > 0 {
		args.ExistedRelativeDirPath = randomItem(allDirs)
	}

	return args
}
