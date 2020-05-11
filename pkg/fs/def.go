package fs

type WorktreeOP int

const (
	Create = WorktreeOP(iota)
	Delete
	Rename
	Override
	TotalFSOP
)

func (op WorktreeOP) String() string {
	return []string{
		"create", "delete", "rename", "override",
	}[op]
}

type WorktreeObject int

const (
	File = WorktreeObject(iota)
	Dir
	TotalFSObject
)

func (o WorktreeObject) String() string {
	return []string{
		"file", "dir",
	}[o]
}

type WorktreeOPArgs struct {
	NewRelativeFilePath     string
	ExistedRelativeFilePath string
	NewRelativeDirPath      string
	ExistedRelativeDirPath  string
	Content                 string
	Offset                  int64
	Size                    int64
}

type Worktree interface {
	AllDirs() []string
	AllFiles() []string
	FileSize(relativePath string) int64
	Apply(ob WorktreeObject, op WorktreeOP, args *WorktreeOPArgs)
}
