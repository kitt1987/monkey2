package cmd

import (
	"bufio"
	"github.com/git-roll/monkey2/pkg/notify"
	"io"
	"os"
	"os/exec"
	"strings"
)

func NewSeq(seqFile, worktree string) *Seq {
	f, err := os.Open(seqFile)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	reader := bufio.NewReader(f)

	var cmds []*Command
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
			c := parseCommand(strings.Join(append(partialLine, string(line)), ""))
			if c != nil {
				cmds = append(cmds, c)
			}

			partialLine = nil
			continue
		}

		c := parseCommand(string(line))
		if c != nil {
			cmds = append(cmds, c)
		}
	}

	return &Seq{CMD: cmds, dir: worktree }
}

func parseCommand(line string) *Command {
	args := strings.Split(line, " ")
	if len(args) == 0 {
		return nil
	}

	return &Command{
		Name: args[0],
		Args: args[1:],
	}
}

type Command struct {
	Name string
	Args []string
}

type Seq struct {
	CMD []*Command
	dir string
}

func (s *Seq) Apply(id int) {
    if id >= len(s.CMD) {
        panic(id)
    }

    c := exec.Command(s.CMD[id].Name, s.CMD[id].Args...)
    c.Env = os.Environ()
    c.Dir = s.dir
    c.Stdout = os.Stdout
    c.Stderr = os.Stderr

    notify.Printf(`💻 Exec command :"%s"`, c.String())
    c.Run()
}
