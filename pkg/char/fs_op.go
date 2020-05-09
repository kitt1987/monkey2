package char

type FileSystemOP int

const (
	FSCreate   = iota
	FSDelete
	FSRename
	FSOverride
	TotalFSOP
)

func (op FileSystemOP) String() string {
	return []string{
		"create", "delete", "rename", "override",
	}[op]
}

type FileSystemObject int

const (
	FSFile = iota
	FSDir
	TotalFSObject
)

func (o FileSystemObject) String() string {
	return []string{
		"file", "dir",
	}[o]
}
