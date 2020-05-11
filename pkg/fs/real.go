package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type real struct {
	baseDir string
}

func (w real) readDir() (dirs, files []string) {
	parents := []string{""}
	for _, parent := range parents {
		path := w.completePath(parent)
		fis, err := ioutil.ReadDir(path)
		if err != nil {
			w.panic(parent, err)
		}

		for _, fi := range fis {
			if fi.IsDir() {
				parents = append(parents, filepath.Join(parent, fi.Name()))
			} else {
				files = append(files, filepath.Join(parent, fi.Name()))
			}
		}
	}

	dirs = parents[1:]
	return
}

func (w real) size(relativePath string) int64 {
	path := w.completePath(relativePath)
	fi, err := os.Lstat(path)
	if err != nil {
		w.panic(path, err)
	}

	return fi.Size()
}

func (w real) createFile(name, text string) {
	path := w.completePath(name)
	if err := ioutil.WriteFile(path, []byte(text), 0755); err != nil {
		w.panic(path, err)
	}
}

func (w real) overrideFile(name, text string, off, size int64) {
	path := w.completePath(name)
	f, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		w.panic(path, err)
	}

	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		w.panic(path, err)
	}

	overriddenLen := fi.Size() - off
	if overriddenLen < 0 {
		w.panic(path, fmt.Errorf("size: %d, offset: %d", fi.Size(), off))
	}

	var overriddenBuf []byte
	if overriddenLen > 0 {
		overriddenBuf = make([]byte, overriddenLen)
		offset := off
		var n int64
		var err error
		for n < overriddenLen && err != io.EOF {
			var m int
			m, err = f.ReadAt(overriddenBuf[n:], offset)
			if m == 0 {
				w.panic(path, err)
			}
			n += int64(m)
			offset += int64(m)
		}
	}

	if len(text) > 0 {
		buf := []byte(text)
		offset := off
		var n int64
		var err error
		for n < int64(len(text)) {
			var m int
			m, err = f.WriteAt(buf[n:], offset)
			if m == 0 {
				w.panic(path, err)
			}

			n += int64(m)
			offset += int64(m)
		}
	}

	if int64(len(overriddenBuf)) > size {
		buf := overriddenBuf[size:]
		var n int64
		var err error
		for n < int64(len(buf)) {
			var m int
			m, err = f.Write(buf[n:])
			if m == 0 {
				w.panic(path, err)
			}

			n += int64(m)
		}
	}
}

func (w real) makeDir(name string) {
	path := w.completePath(name)
	if err := os.MkdirAll(path, 0755); err != nil {
		w.panic(path, err)
	}
}

func (w real) delete(name string) {
	path := w.completePath(name)
	if err := os.RemoveAll(path); err != nil {
		w.panic(path, err)
	}
}

func (w real) rename(origin, target string) {
	originPath := w.completePath(origin)
	targetPath := w.completePath(target)
	if err := os.Rename(originPath, targetPath); err != nil {
		w.panic(originPath, err)
	}
}

func (w real) readFile(name string) string {
	path := w.completePath(name)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		w.panic(path, err)
	}
	return string(bytes)
}

func (w real) completePath(name string) (path string) {
	return filepath.Join(w.baseDir, name)
}

func (w real) panic(path string, err error) {
	panic(fmt.Sprintf("%s:%s", path, err))
}
