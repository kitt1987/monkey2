package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func NewWorktree(workDir string) Worktree {
	fi, err := os.Lstat(workDir)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(fmt.Sprintf("%s:%s", workDir, err))
		}

		err = os.MkdirAll(workDir, 0755)
		if err != nil {
			panic(fmt.Sprintf("%s:%s", workDir, err))
		}
	} else if !fi.IsDir() {
		panic(fmt.Sprintf("%s:not a directory", workDir))
	}

	return &worktree{baseDir: workDir}
}

type worktree struct {
	baseDir string
}

func (w worktree) readDir() (dirs, files []string) {
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

func (w worktree) AllDirs() []string {
	dirs, _ := w.readDir()
	return dirs
}

func (w worktree) AllFiles() []string {
	_, files := w.readDir()
	return files
}

func (w worktree) FileSize(relativePath string) int64 {
	path := w.completePath(relativePath)
	fi, err := os.Lstat(path)
	if err != nil {
		w.panic(path, err)
	}

	return fi.Size()
}

func (w worktree) Apply(ob WorktreeObject, op WorktreeOP, args *WorktreeOPArgs) {
	switch ob {
	case File:
		w.applyFile(op, args)
	case Dir:
		w.applyDir(op, args)
	default:
		panic(ob)
	}
}

func (w worktree) applyFile(op WorktreeOP, args *WorktreeOPArgs) {
	switch op {
	case Create:
		fmt.Printf(`💻 Create file "%s" with content:
+++
%s
+++
`, args.NewRelativeFilePath, args.Content)

		w.createFile(args.NewRelativeFilePath, args.Content)
	case Delete:
		fmt.Printf(`💻 Unlink "%s"`+"\n", args.ExistedRelativeFilePath)
		w.delete(args.ExistedRelativeFilePath)
	case Rename:
		fmt.Printf(`💻️ Rename file "%s" to "%s"`+"\n", args.ExistedRelativeFilePath, args.NewRelativeFilePath)
		w.rename(args.ExistedRelativeFilePath, args.NewRelativeFilePath)
	case Override:
		fmt.Printf(`💻️ Overwrite file "%s", replace %d bytes content from byte %d with:
+++
%s
+++
`,
			args.ExistedRelativeFilePath, args.Size, args.Offset, args.Content)
		w.overrideFile(args.ExistedRelativeFilePath, args.Content, args.Offset, args.Size)
	default:
		panic(op)
	}
}

func (w worktree) applyDir(op WorktreeOP, args *WorktreeOPArgs) {
	switch op {
	case Create:
		fmt.Printf(`💻 Mkdir "%s"`+"\n", args.NewRelativeDirPath)
		w.makeDir(args.NewRelativeDirPath)
	case Delete:
		fmt.Printf(`💻 Unlink "%s"`+"\n", args.ExistedRelativeDirPath)
		w.delete(args.ExistedRelativeDirPath)
	case Rename:
		fmt.Printf(`💻 Rename dir "%s" to "%s"`+"\n", args.ExistedRelativeDirPath, args.NewRelativeDirPath)
		w.rename(args.ExistedRelativeDirPath, args.NewRelativeDirPath)
	default:
		panic(op)
	}
}

func (w worktree) createFile(name, text string) {
	path := w.completePath(name)
	if err := ioutil.WriteFile(path, []byte(text), 0755); err != nil {
		w.panic(path, err)
	}
}

func (w worktree) overrideFile(name, text string, off, size int64) {
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

func (w worktree) makeDir(name string) {
	path := w.completePath(name)
	if err := os.MkdirAll(path, 0755); err != nil {
		w.panic(path, err)
	}
}

func (w worktree) delete(name string) {
	path := w.completePath(name)
	if err := os.RemoveAll(path); err != nil {
		w.panic(path, err)
	}
}

func (w worktree) rename(origin, target string) {
	originPath := w.completePath(origin)
	targetPath := w.completePath(target)
	if err := os.Rename(originPath, targetPath); err != nil {
		w.panic(originPath, err)
	}
}

func (w worktree) completePath(name string) (path string) {
	return filepath.Join(w.baseDir, name)
}

func (w worktree) panic(path string, err error) {
	panic(fmt.Sprintf("%s:%s", path, err))
}
