package background

import (
	"os"
	"os/signal"
	"syscall"
)

// Wait blocked and waits signal from os
func Wait() os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	s := <-quit
	return s
}
