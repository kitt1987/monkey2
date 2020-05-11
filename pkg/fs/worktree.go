package fs

import "fmt"

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

		w.under.createFile(args.NewRelativeFilePath, args.Content)
	case Delete:
		fmt.Printf(`💻 Unlink "%s"`+"\n", args.ExistedRelativeFilePath)
		w.under.delete(args.ExistedRelativeFilePath)
	case Rename:
		fmt.Printf(`💻️ Rename file "%s" to "%s"`+"\n", args.ExistedRelativeFilePath, args.NewRelativeFilePath)
		w.under.rename(args.ExistedRelativeFilePath, args.NewRelativeFilePath)
	case Override:
		fmt.Printf(`💻️ Overwrite file "%s", replace %d bytes content from byte %d with:
+++
%s
+++
`,
			args.ExistedRelativeFilePath, args.Size, args.Offset, args.Content)
		w.under.overrideFile(args.ExistedRelativeFilePath, args.Content, args.Offset, args.Size)
	default:
		panic(op)
	}
}

func (w worktree) applyDir(op WorktreeOP, args *WorktreeOPArgs) {
	switch op {
	case Create:
		fmt.Printf(`💻 Mkdir "%s"`+"\n", args.NewRelativeDirPath)
		w.under.makeDir(args.NewRelativeDirPath)
	case Delete:
		fmt.Printf(`💻 Unlink "%s"`+"\n", args.ExistedRelativeDirPath)
		w.under.delete(args.ExistedRelativeDirPath)
	case Rename:
		fmt.Printf(`💻 Rename dir "%s" to "%s"`+"\n", args.ExistedRelativeDirPath, args.NewRelativeDirPath)
		w.under.rename(args.ExistedRelativeDirPath, args.NewRelativeDirPath)
	default:
		panic(op)
	}
}

func (w worktree) validateFSStructure() {
	if w.mirror == nil {
		return
	}

	mirrorDirs, mirrorFiles := w.mirror.readDir()
	dirs, files := w.under.readDir()
	if !equalStringSlices(dirs, mirrorDirs) {
		panic(fmt.Sprintf("mirror:%#v \n real:%#v", mirrorDirs, dirs))
	}

	if !equalStringSlices(files, mirrorFiles) {
		panic(fmt.Sprintf("mirror:%#v \n real:%#v", mirrorFiles, files))
	}
}

func (w worktree) validateFile(name string) {
	if w.mirror == nil {
		return
	}

	mirrorContent := w.mirror.readFile(name)
	content := w.under.readFile(name)
	if content != mirrorContent {
		panic(fmt.Sprintf("file: %s \n mirror:%s, \n real:%s", name, mirrorContent, content))
	}
}

func equalStringSlices(a, b []string) bool {

}
