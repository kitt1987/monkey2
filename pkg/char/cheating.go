package char

import (
    "github.com/git-roll/monkey2/pkg/notify"
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

func bareClone(repo, local string) (commits[]string, err error) {
    notify.Printf("🚁 Clone target repo %s\n", repo)
    cmd := exec.Command("git", "clone", "--bare", repo, local)
    cmd.Stdout = notify.Writer()
    cmd.Stderr = notify.Writer()
    if err = cmd.Run(); err != nil {
        return
    }

    cmdListCommits := exec.Command("git", "--git-dir", local, "")
}

type cheatingMonkey struct {
    worktree string
    targetLocal string
}

func (c *cheatingMonkey) Work() {
    cmd := exec.Command("git",
        "--work-tree", c.worktree, "--git-dir", c.targetLocal,
        "checkout", )
    cmd.Stdout = notify.Writer()
    cmd.Stderr = notify.Writer()
    return cmd.Run()
}
