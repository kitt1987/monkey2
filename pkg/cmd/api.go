package cmd

import (
    "bufio"
    "io"
    "os"
    "strings"
)

func NewSeq(seqFile string) (*Seq, error) {
    f, err := os.Open(seqFile)
    if err != nil {
        return nil, err
    }

    defer f.Close()
    reader := bufio.NewReader(f)

    var cmds []Command
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
            cmds = append(cmds, parseCommand(strings.Join(append(partialLine, string(line)), "")))
            partialLine = nil
            continue
        }

        cmds = append(cmds, parseCommand(string(line)))
    }

    return &Seq{cmds: cmds}, nil
}

func parseCommand(line string) Command {
    
}

type Command struct {
    Name string
    Args []string
}

type Seq struct {
    cmds []Command
}
