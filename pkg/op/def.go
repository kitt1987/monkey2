package op

type WorktreeOP int

const (
	FSCreate = iota
	FSDelete
	FSRename
	FSOverride
	TotalFSOP
)

func (op WorktreeOP) String() string {
	return []string{
		"create", "delete", "rename", "override",
	}[op]
}

type WorktreeObject int

const (
	FSFile = iota
	FSDir
	TotalFSObject
)

func (o WorktreeObject) String() string {
	return []string{
		"file", "dir",
	}[o]
}

type WorktreeOPArgs struct {
	RelativeFilePath    string
	RelativeDirPath     string
	RelativeRenamedFile string
	RelativeRenamedDir  string
	Content             string
	Offset              int64
	Size                int64
}
