package helpers

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/team-scaletech/common/logging"
)

func GracefulStop(callback func(context.Context) error) {
	zlog := logging.GetLog()

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGTERM)
	<-gracefulStop

	zlog.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := callback(ctx); err != nil {
		zlog.Fatal().Err(err).Msgf("Server forced to shutdown: %+v", err)
	}

	zlog.Info().Msg("Server exiting")
}
