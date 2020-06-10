package ws

import (
	"github.com/git-roll/monkey2/pkg/notify"
	"github.com/gorilla/websocket"
	"net/http"
)

type websocketSession struct {
	conn   *websocket.Conn
	closed chan struct{}
}

type websocketWriter struct {
	sessions []websocketSession
	recv     chan []byte
}

func (wsw *websocketWriter) Write(p []byte) (n int, err error) {
	wsw.recv <- p
	return len(p), err
}

func (wsw *websocketWriter) Start() {
	for {
		bytes, ok := <-wsw.recv
		if !ok {
			return
		}

		liveSessions := make([]websocketSession, 0, len(wsw.sessions))
		for i := range wsw.sessions {
			s := wsw.sessions[i]
			if err := s.conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
				close(s.closed)
				continue
			}

			liveSessions = append(liveSessions, s)
		}

		wsw.sessions = liveSessions
	}
}

func (wsw *websocketWriter) Close() error {
	close(wsw.recv)
	for i := range wsw.sessions {
		s := &wsw.sessions[i]
		close(s.closed)
	}

	return nil
}

type websocketHandle struct {
	writer *websocketWriter
	*notify.Mirror
}

func newWSWriter() *websocketHandle {
	w := &websocketWriter{
		recv: make(chan []byte, 64),
	}

	go w.Start()
	return &websocketHandle{
		writer: w,
		Mirror: notify.NewMirrorNotifier(w),
	}
}

var (
	upgrade = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (wsw *websocketHandle) addConn(conn *websocket.Conn) <-chan struct{} {
	closed := make(chan struct{})
	wsw.writer.sessions = append(wsw.writer.sessions, websocketSession{
		conn:   conn,
		closed: closed,
	})

	return closed
}

func (wsw *websocketHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	defer conn.Close()

	closed := wsw.addConn(conn)
	<-closed
}
