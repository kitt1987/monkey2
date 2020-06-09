package notify

import "strings"

func NewMirrorNotifier(notifier io.WriteCloser) *Mirror {
	return &Mirror{
		notifier:     notifier,
		lastNotes:    make([]string, 0, 50),
		maxLastNotes: 50,
	}
}

type Mirror struct {
	notifier     io.WriteCloser
	lastNotes    []string
	maxLastNotes int
}

func (n2 *Mirror) Write(p []byte) (n int, err error) {
	if len(n2.lastNotes) < n2.maxLastNotes {
		n2.lastNotes = append(n2.lastNotes, string(p))
	} else {
		n2.lastNotes[0] = ""
		n2.lastNotes = append(n2.lastNotes[1:], string(p))
	}

	return n2.notifier.Write(p)
}

func (n2 *Mirror) Close() error {
	return n2.notifier.Close()
}

func (n2 Mirror) LastNotes() string {
	return strings.Join(n2.lastNotes, "\n")
}
