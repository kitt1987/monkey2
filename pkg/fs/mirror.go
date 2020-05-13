package fs

import "fmt"

type mirrorElem struct {
	dir     bool
	content []byte
}

func newMirror(initDirs, initFiles []string, initFileContents map[string]string) (m memMirror) {
	m = make(memMirror)
	for _, dir := range initDirs {
		m[dir] = &mirrorElem{
			dir: true,
		}
	}

	for _, file := range initFiles {
		m[file] = &mirrorElem{
			content: []byte(initFileContents[file]),
		}
	}

	return m
}

type memMirror map[string]*mirrorElem

func (m memMirror) readDir() (dirs, files []string) {
	dirs = make([]string, 0, len(m))
	for p, e := range m {
		if e.dir {
			dirs = append(dirs, p)
		} else {
			files = append(files, p)
		}
	}

	return
}

func (m memMirror) size(relativePath string) int64 {
	elem, found := m[relativePath]
	if !found {
		panic(fmt.Sprintf(`%s not found`, relativePath))
	}

	if elem.dir {
		panic(fmt.Sprintf(`%s is a directory`, relativePath))
	}

	return int64(len(elem.content))
}

func (m memMirror) createFile(name, text string) {
	if _, found := m[name]; found {
		panic(fmt.Sprintf(`%s already exists`, name))
	}

	m[name] = &mirrorElem{
		content: []byte(text),
	}
}

func (m memMirror) overrideFile(name, text string, off, size int64) {
	elem, found := m[name]
	if !found {
		panic(fmt.Sprintf(`%s not found`, name))
	}

	if elem.dir {
		panic(fmt.Sprintf(`%s is a directory`, name))
	}

	var tail []byte
	if (off + size) < int64(len(elem.content)) {
		tail = elem.content[off+size:]
	}

	elem.content = append(elem.content[:off], append([]byte(text), tail...)...)
}

func (m memMirror) makeDir(name string) {
	if _, found := m[name]; found {
		panic(fmt.Sprintf(`%s already exists`, name))
	}

	m[name] = &mirrorElem{
		dir: true,
	}
}

func (m memMirror) delete(name string) {
	if _, found := m[name]; !found {
		panic(fmt.Sprintf(`%s not found`, name))
	}

	delete(m, name)
}

func (m memMirror) rename(origin, target string) {
	elem, found := m[origin]
	if !found {
		panic(fmt.Sprintf(`%s not found`, origin))
	}

	if _, found := m[target]; found {
		panic(fmt.Sprintf(`%s already exists`, target))
	}

	delete(m, origin)
	m[target] = elem
}

func (m memMirror) readFile(name string) string {
	elem, found := m[name]
	if !found {
		panic(fmt.Sprintf(`%s not found`, name))
	}

	if elem.dir {
		panic(fmt.Sprintf(`%s is a directory`, name))
	}

	return string(elem.content)
}
