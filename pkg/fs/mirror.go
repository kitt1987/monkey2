package fs

import "fmt"

type mirrorElem struct {
	dir     bool
	content []byte
}

type wtMirror map[string]mirrorElem

func (m wtMirror) AllDirs() (dirs []string) {
	dirs = make([]string, 0, len(m))
	for p, e := range m {
		if e.dir {
			dirs = append(dirs, p)
		}
	}

	return
}

func (m wtMirror) AllFiles() (files []string) {
	files = make([]string, 0, len(m))
	for p, e := range m {
		if !e.dir {
			files = append(files, p)
		}
	}

	return
}

func (m wtMirror) FileSize(relativePath string) int64 {
	elem, found := m[relativePath]
	if !found {
		panic(fmt.Sprintf(`%s not found`, relativePath))
	}

	if elem.dir {
		panic(fmt.Sprintf(`%s is a directory`, relativePath))
	}

	return int64(len(elem.content))
}

func (m wtMirror) createFile(name, text string) {
	if _, found := m[name]; found {
		panic(fmt.Sprintf(`%s already exists`, name))
	}

	m[name] = mirrorElem{
		content: []byte(text),
	}
}

func (m wtMirror) overrideFile(name, text string, off, size int64) {
	elem, found := m[name]
	if !found {
		panic(fmt.Sprintf(`%s not found`, name))
	}

	if elem.dir {
		panic(fmt.Sprintf(`%s is a directory`, name))
	}

	tail := elem.content[:off+size]
	elem.content = append(elem.content[:off], []byte(text)...)
	elem.content = append(elem.content, tail...)
}

func (m wtMirror) makeDir(name string) {
	if _, found := m[name]; found {
		panic(fmt.Sprintf(`%s already exists`, name))
	}

	m[name] = mirrorElem{
		dir: true,
	}
}

func (m wtMirror) delete(name string) {
	if _, found := m[name]; !found {
		panic(fmt.Sprintf(`%s not found`, name))
	}

	delete(m, name)
}

func (m wtMirror) rename(origin, target string) {
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
