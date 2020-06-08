package fs

import (
	"fmt"
	"github.com/git-roll/monkey2/pkg/notify"
)

type worktree struct {
	under  underneath
	mirror underneath
}

func (w worktree) AllDirs() []string {
	dirs, _ := w.under.readDir()
	return dirs
}

func (w worktree) AllFiles() []string {
	_, files := w.under.readDir()
	return files
}

func (w worktree) FileSize(relativePath string) int64 {
	return w.under.size(relativePath)
}

func (w worktree) Apply(ob WorktreeObject, op WorktreeOP, args *WorktreeOPArgs) {
	if err := w.validateFSStructure(); err != nil {
		return
	}

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
		notify.Printf(`💻 Create file "%s" with content:
+++
%s
+++
`, args.NewRelativeFilePath, args.Content)

		if w.mirror != nil {
			w.mirror.createFile(args.NewRelativeFilePath, args.Content)
		}

		w.under.createFile(args.NewRelativeFilePath, args.Content)
	case Delete:
		notify.Printf(`💻 Unlink "%s"`+"\n", args.ExistedRelativeFilePath)

		if w.mirror != nil {
			w.mirror.delete(args.ExistedRelativeFilePath)
		}

		w.under.delete(args.ExistedRelativeFilePath)
	case Rename:
		notify.Printf(`💻️ Rename file "%s" to "%s"`+"\n", args.ExistedRelativeFilePath, args.NewRelativeFilePath)

		if w.mirror != nil {
			w.mirror.rename(args.ExistedRelativeFilePath, args.NewRelativeFilePath)
		}

		w.under.rename(args.ExistedRelativeFilePath, args.NewRelativeFilePath)
	case Override:
		notify.Printf(`💻️ Overwrite file "%s", replace %d bytes content from byte %d with:
+++
%s
+++
`,
			args.ExistedRelativeFilePath, args.Size, args.Offset, args.Content)

		w.validateFile(args.ExistedRelativeFilePath)

		if w.mirror != nil {
			w.mirror.overrideFile(args.ExistedRelativeFilePath, args.Content, args.Offset, args.Size)
		}

		w.under.overrideFile(args.ExistedRelativeFilePath, args.Content, args.Offset, args.Size)
	default:
		panic(op)
	}
}

func (w worktree) applyDir(op WorktreeOP, args *WorktreeOPArgs) {
	switch op {
	case Create:
		notify.Printf(`💻 Mkdir "%s"`+"\n", args.NewRelativeDirPath)

		if w.mirror != nil {
			w.mirror.makeDir(args.NewRelativeDirPath)
		}

		w.under.makeDir(args.NewRelativeDirPath)
	case Delete:
		notify.Printf(`💻 Unlink "%s"`+"\n", args.ExistedRelativeDirPath)

		if w.mirror != nil {
			w.mirror.delete(args.ExistedRelativeDirPath)
		}

		w.under.delete(args.ExistedRelativeDirPath)
	case Rename:
		notify.Printf(`💻 Rename dir "%s" to "%s"`+"\n", args.ExistedRelativeDirPath, args.NewRelativeDirPath)

		if w.mirror != nil {
			w.mirror.rename(args.ExistedRelativeDirPath, args.NewRelativeDirPath)
		}

		w.under.rename(args.ExistedRelativeDirPath, args.NewRelativeDirPath)
	default:
		panic(op)
	}
}

func (w worktree) validateFSStructure() error {
	if w.mirror == nil {
		return nil
	}

	mirrorDirs, mirrorFiles := w.mirror.readDir()
	dirs, files := w.under.readDir()
	if err := equalStringSlices(dirs, mirrorDirs); err != nil {
		panic(fmt.Sprintf("err: %s\n(a)real:%#v\n(b)mirror:%#v\n", err, dirs, mirrorDirs))
	}

	if err := equalStringSlices(files, mirrorFiles); err != nil {
		panic(fmt.Sprintf("err: %s\n(a)real:%#v\n(b)mirror:%#v\n", err, files, mirrorFiles))
	}

	return nil
}

func (w worktree) validateFile(name string) {
	if w.mirror == nil {
		return
	}

	mirrorContent := w.mirror.readFile(name)
	content := w.under.readFile(name)
	if content != mirrorContent {
		panic(fmt.Sprintf("file: %s\n---\nmirror:%s\n---\nreal:%s\n", name, mirrorContent, content))
	}
}

func equalStringSlices(a, b []string) error {
	if len(a) != len(b) {
		return fmt.Errorf("#a=%d, #b=%d", len(a), len(b))
	}

	mapA := make(map[string]bool, len(a))
	for _, ai := range a {
		if mapA[ai] {
			return fmt.Errorf(`slice a has a duplicate value "%s"`, ai)
		}

		mapA[ai] = true
	}

	for _, bi := range b {
		if !mapA[bi] {
			return fmt.Errorf(`"%s" didn't found in slice a'`, bi)
		}

		delete(mapA, bi)
	}

	return nil
}
