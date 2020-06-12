package char

import (
    "bufio"
    "github.com/git-roll/monkey2/pkg/cmd"
    "github.com/git-roll/monkey2/pkg/conf"
    "github.com/git-roll/monkey2/pkg/notify"
    "io"
    "os/exec"
    "path/filepath"
    "strings"
)

func Cheating(worktree, targetRepo string, recover func(string)) Monkey {
	localRepo := getLocalRepoPath(worktree)
	commits, err := bareClone(targetRepo, localRepo)
	if err != nil {
		panic(err)
	}

	m := &cheatingMonkey{
        targetLocal: localRepo,
        worktree:    worktree,
        commits:     commits,
    }

    pr := conf.CmdBuildPR()
    if len(pr) > 0 {
        m.buildPR = cmd.NewSeqFromText(pr, worktree)
    }

	return &monkey{
        recover: recover,
        monkeyChar: m,
    }
}

func getLocalRepoPath(worktree string) string {
	return filepath.Join(filepath.Base(worktree), "target.git")
}

type cheatedCommit struct {
	Hash   string
	Merged bool
}

func bareClone(repo, local string) (commits []*cheatedCommit, err error) {
	notify.Printf("🚁 Clone target repo %s\n", repo)
	cmd := exec.Command("git", "clone", "--bare", repo, local)
	cmd.Stdout = notify.Writer()
	cmd.Stderr = notify.Writer()
	if err = cmd.Run(); err != nil {
		return
	}

	cmdListCommits := exec.Command("git", "--git-dir", local, "log", "--reverse", "--format=%h %p")
	out, err := cmdListCommits.Output()
	if err != nil {
		return
	}

	reader := bufio.NewReader(strings.NewReader(string(out)))

	var partialLine []string
	for {
		line, remaining, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		if remaining {
			partialLine = append(partialLine, string(line))
			continue
		}

		var wholeLine string
		if len(partialLine) > 0 {
			wholeLine = strings.Join(append(partialLine, string(line)), "")
			partialLine = nil
		} else {
			wholeLine = string(line)
		}

		segs := strings.Split(wholeLine, " ")
		commits = append(commits, &cheatedCommit{
            Hash:   segs[0],
            Merged: len(segs) > 2,
        })
	}

	return
}

type cheatingMonkey struct {
	worktree    string
	targetLocal string
	commits []*cheatedCommit
    buildPR *cmd.Seq
}

func (c *cheatingMonkey) Work() {
	if len(c.commits) == 0 {
		notify.Printf("🚁 All code is cheated")
		return
	}

	notify.Printf("👻 Cheat commit %s", c.commits[0])
	cmd := exec.Command("git",
		"--work-tree", c.worktree, "--git-dir", c.targetLocal,
		"checkout", c.commits[0].Hash)
	cmd.Stdout = notify.Writer()
	cmd.Stderr = notify.Writer()
	err := cmd.Run()
	if err != nil {
		panic(c.commits[0])
	}

	if c.commits[0].Merged && c.buildPR != nil {
        c.buildPR.Apply(0)
    }

	c.commits[0] = nil
	c.commits = c.commits[1:]
}
