package cmd

import (
	"context"

	"awesomeProjectCr/cmd/app"

	"github.com/rs/zerolog/log"
)

// StartServer : Starts the server both the gRPC and REST servers.
func StartServer() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	startServer(ctx, cancel)

	<-ctx.Done()
	log.Info().Msg("All servers stopped")
}

func startServer(ctx context.Context, cancel context.CancelFunc) {
	s := app.New()

	s.Start(ctx, cancel)
}
