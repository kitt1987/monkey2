package op

type fsElem struct {
	dir bool
	content string
}

type fsMirror map[string]fsElem

func (m fsMirror) AllDirs() []string {
	
}

func (m fsMirror) AllFiles() []string {

}

func (m fsMirror) FileSize(relativePath string) int64 {

}

func (m fsMirror) createFile(name, text string) {
	
}

func (m fsMirror) overrideFile(name, text string, off, size int64) {
	
}

func (m fsMirror) makeDir(name string) {
	
}

func (m fsMirror) delete(name string) {
	
}

func (m fsMirror) rename(origin, target string) {
	
}
