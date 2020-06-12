package notify

import (
	"fmt"
	"io"
	"os"
)

var no = io.Writer(os.Stdout)

func Set(notifier io.Writer) {
	no = notifier
}

func Printf(text string, a ...interface{}) {
	fmt.Fprintf(no, text, a...)
	fmt.Printf(text, a...)
}

func Writer() io.Writer {
	return no
}
