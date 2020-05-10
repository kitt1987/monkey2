package op

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Worktree struct {
	baseDir string
}

func (w Worktree) readDir() (dirs, files []string) {
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

func (w Worktree) AllDirs() []string {
	dirs, _ := w.readDir()
	return dirs
}

func (w Worktree) AllFiles() []string {
	_, files := w.readDir()
	return files
}

func (w Worktree) FileSize(relativePath string) int64 {
	path := w.completePath(relativePath)
	fi, err := os.Lstat(path)
	if err != nil {
		w.panic(path, err)
	}

	return fi.Size()
}

func (w Worktree) Apply(ob WorktreeObject, op WorktreeOP, args *WorktreeOPArgs) {
	switch ob {
	case FSFile:
		w.applyFile(op, args)
	case FSDir:
		w.applyDir(op, args)
	default:
		panic(ob)
	}
}

func (w Worktree) applyFile(op WorktreeOP, args *WorktreeOPArgs) {
	switch op {
	case FSCreate:
		w.createFile(args.NewRelativeFilePath, args.Content)
	case FSDelete:
		w.delete(args.ExistedRelativeFilePath)
	case FSRename:
		w.rename(args.ExistedRelativeFilePath, args.NewRelativeFilePath)
	case FSOverride:
		w.overrideFile(args.ExistedRelativeFilePath, args.Content, args.Offset, args.Size)
	default:
		panic(op)
	}
}

func (w Worktree) applyDir(op WorktreeOP, args *WorktreeOPArgs) {
	switch op {
	case FSCreate:
		w.makeDir(args.NewRelativeDirPath)
	case FSDelete:
		w.delete(args.ExistedRelativeDirPath)
	case FSRename:
		w.rename(args.ExistedRelativeDirPath, args.NewRelativeDirPath)
	default:
		panic(op)
	}
}

func (w Worktree) createFile(name, text string) {
	path := w.completePath(name)
	if err := ioutil.WriteFile(path, []byte(text), 0755); err != nil {
		w.panic(path, err)
	}
}

func (w Worktree) overrideFile(name, text string, off, size int64) {
	path := w.completePath(name)
	f, err := os.OpenFile(path, os.O_RDWR, 0755)
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
		overriddenBuf = make([]byte, 0, overriddenLen)
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

func (w Worktree) makeDir(name string) {
	path := w.completePath(name)
	if err := os.MkdirAll(path, 0755); err != nil {
		w.panic(path, err)
	}
}

func (w Worktree) delete(name string) {
	path := w.completePath(name)
	if err := os.RemoveAll(path); err != nil {
		w.panic(path, err)
	}
}

func (w Worktree) rename(origin, target string) {
	originPath := w.completePath(origin)
	targetPath := w.completePath(target)
	if err := os.Rename(originPath, targetPath); err != nil {
		w.panic(originPath, err)
	}
}

func (w Worktree) completePath(name string) (path string) {
	path = filepath.Join(w.baseDir, name)
	if _, err := os.Lstat(path); !os.IsNotExist(err) {
		panic(path)
	}

	return
}

func (w Worktree) panic(path string, err error) {
	panic(fmt.Sprintf("%s:%s", path, err))
}