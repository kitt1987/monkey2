package char

import (
    "os/exec"
    "path/filepath"
)

func Cheating(worktree, targetRepo string, recover func(string)) Monkey {
    localRepo := getLocalRepoPath(worktree)
    if err := bareClone(targetRepo, localRepo); err != nil {
        panic(err)
    }

    return &monkey{
        recover:    recover,
        monkeyChar: &cheatingMonkey{
            targetLocal: localRepo,
            worktree: worktree,
        },
    }
}

func getLocalRepoPath(worktree string) string {
    return filepath.Join(filepath.Base(worktree), "target.git")
}

func bareClone(repo, local string) error {
    cmd := exec.Command("git", "clone", "--bare", repo, local)
}

type cheatingMonkey struct {
    worktree string
    targetLocal string
}

func (c *cheatingMonkey) Work() {
    panic("implement me")
}
