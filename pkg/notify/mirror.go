package notify

import (
    "io"
    "strings"
)

func NewMirrorNotifier(notifier io.WriteCloser) *Mirror {
	return &Mirror{
		notifier:     notifier,
		latestNotes:  make([]string, 0, 50),
		maxLastNotes: 50,
	}
}

type Mirror struct {
	notifier     io.WriteCloser
	latestNotes  []string
	maxLastNotes int
}

func (n2 *Mirror) Write(p []byte) (n int, err error) {
	if len(n2.latestNotes) < n2.maxLastNotes {
		n2.latestNotes = append(n2.latestNotes, string(p))
	} else {
		n2.latestNotes[0] = ""
		n2.latestNotes = append(n2.latestNotes[1:], string(p))
	}

	return n2.notifier.Write(p)
}

func (n2 *Mirror) Close() error {
	return n2.notifier.Close()
}

func (n2 Mirror) JoinNotices() string {
	return strings.Join(n2.latestNotes, "\n")
}

func (n2 Mirror) Notices() []string {
    return n2.latestNotes
}
