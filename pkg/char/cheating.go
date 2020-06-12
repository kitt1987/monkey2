package char

import (
	"bufio"
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

	return &monkey{
		recover: recover,
		monkeyChar: &cheatingMonkey{
			targetLocal: localRepo,
			worktree:    worktree,
			commits:     commits,
		},
	}
}

func getLocalRepoPath(worktree string) string {
	return filepath.Join(filepath.Base(worktree), "target.git")
}

func bareClone(repo, local string) (commits []string, err error) {
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

		if len(partialLine) > 0 {
			commits = append(commits, strings.Join(append(partialLine, string(line)), ""))
			partialLine = nil
			continue
		}

		commits = append(commits, string(line))
	}

	return
}

type cheatingMonkey struct {
	worktree    string
	targetLocal string
	commits     []string
}

func (c *cheatingMonkey) Work() {
    if len(c.commits) == 0 {
        notify.Printf("🚁 All code is cheated")
        return
    }

    notify.Printf("👻 Cheat commit %s", c.commits[0])
	cmd := exec.Command("git",
		"--work-tree", c.worktree, "--git-dir", c.targetLocal,
		"checkout", c.commits[0])
	cmd.Stdout = notify.Writer()
	cmd.Stderr = notify.Writer()
	err := cmd.Run()
	if err != nil {
	    panic(c.commits[0])
    }

    c.commits[0] = ""
    c.commits = c.commits[1:]
}
