package cmd

import (
	"context"
	"os"

	"awesomeProjectCr/cmd/app"

	"github.com/rs/zerolog/log"
)

func StartServer() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := app.New()
	s.Start(ctx, cancel)

	<-ctx.Done()

	log.Info().Msg("all servers stopped")

	if err := ctx.Err(); err != nil && err != context.Canceled {
		log.Error().Err(err).Msg("server stopped due to unexpected error")
		os.Exit(1)
	}

	os.Exit(0)
}
