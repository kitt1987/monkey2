package op

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Worktree struct {
	baseDir string
}

func (w Worktree) NewFile(name, text string) {
	path := filepath.Join(w.baseDir, name)
	if _, err := os.Lstat(path); !os.IsNotExist(err) {
		panic(path)
	}

	if err := ioutil.WriteFile(path, []byte(text), 0755); err != nil {
		panic(fmt.Sprintf("%s:%s", path, err))
	}
}

func (w Worktree) NewDir(name string) {

}
