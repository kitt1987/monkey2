package ws

import (
	"fmt"
	"github.com/git-roll/monkey2/pkg/conf"
	"github.com/git-roll/monkey2/pkg/util"
	"io"
	"net"
	"net/http"
	"time"
)

func NewServer() *Server {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if _, err := writer.Write([]byte(homePage)); err != nil {
			fmt.Printf("fail to write the client: %s", err)
		}
	})

	s := &Server{}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.WebSocketPort()))
	if err != nil {
		panic(err)
	}

	s.lis = lis
	return s
}

type Server struct {
	lis net.Listener
}

func (s *Server) SidecarNotifier() io.WriteCloser {
	writer := newWSWriter()
	http.Handle("/sidecar", writer)
	return writer
}

func (s *Server) MonkeyNotifier() io.WriteCloser {
	writer := newWSWriter()
	http.Handle("/monkey", writer)
	return writer
}

func (s *Server) Run(stopC <-chan struct{}) {
	done := util.Start(func() {
		http.Serve(s.lis, nil)
	})

	<-stopC
	s.lis.Close()
	<-done
}

var homePage = fmt.Sprintf(`<!doctype html>
<html lang="en">
  <head>
    <title>Monkey</title>
    <meta charset="UTF-8" />
    <meta name="description" content="🐵" />
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.0/css/bootstrap.min.css" integrity="sha384-9aIt2nRpC12Uk9gS9baDl411NQApFmC26EwAOH8WgZl5MYYxFfc+NcPb1dKGj7Sk" crossorigin="anonymous">
	<link href="https://cdn.bootcdn.net/ajax/libs/xterm/3.14.5/xterm.min.css" rel="stylesheet">
	<script src="https://cdn.bootcdn.net/ajax/libs/xterm/3.14.5/xterm.min.js"></script>
  </head>
  <body>
    <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
      <a class="navbar-brand" href="#">Monkey</a>
      <div class="collapse navbar-collapse" id="navbarNav">
        <ul class="navbar-nav mr-auto">
          <li class="nav-item">
            <a class="nav-link" href="#">INSANE(booted at %s)</a>
          </li>
        </ul>
        <span>🐵</span>
      </div>
    </nav>
    <div class="container-fluid">
      <div class="row justify-content-around mt-3">
        <div id="monkey" class="col-12 col-sm-12 col-md-12 col-lg-11 col-xl-8 shadow p-3 m-1 bg-white">
          <p class="lead">🐵 Monkey</p>
          <div id="monkey_console"></div>
        </div>

        <div id="sidecar" class="col-12 col-sm-12 col-md-12 col-lg-11 col-xl-8 shadow p-3 m-1 bg-white">
          <p class="lead">📺 Sidecar</p>
          <div id="sidecar_console"></div>
        </div>
      </div>
    </div>
    <script>
      function connectConsole(console, ws) {
        var monkeyTerm = new Terminal({
          cols: 100,
        });
        monkeyTerm.open(document.getElementById(console));

        monkey = new WebSocket(ws);
        monkey.onmessage = function (e) {
            monkeyTerm.write(e.data.replace(/\n/g, '\n\r'));
        };
      }
      

      connectConsole("monkey_console", "ws://" + window.location.host + "/monkey")
      connectConsole("sidecar_console", "ws://" + window.location.host + "/sidecar")
    </script>
  </body>
</html>
`, time.Now().Format(time.RFC1123Z))
