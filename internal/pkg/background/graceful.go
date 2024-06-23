package background

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
)

func Wait(logger log.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	s := <-quit
	logger.Info().Msgf("service shutting down; received signal %s", s)
}
