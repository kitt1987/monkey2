package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	signCh := make(chan os.Signal, 3)
	signal.Ignore(syscall.SIGPIPE)
	signal.Notify(signCh, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)

	for {
		select {
		case <-signCh:
		}
	}
}
