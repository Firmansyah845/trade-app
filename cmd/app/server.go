package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"awesomeProjectCr/internal/config"

	"github.com/rs/zerolog/log"
)

type Server struct {
	server *http.Server
}

func New() *Server {
	handler := router()

	server := &Server{
		server: &http.Server{
			Addr:         ":" + strconv.Itoa(config.Server.Port),
			Handler:      handler,
			ReadTimeout:  config.Server.ReadTimeout,
			WriteTimeout: config.Server.WriteTimeout,
			IdleTimeout:  config.Server.IdleTimeout,
		},
	}

	return server
}

func (s *Server) Start(ctx context.Context, cancel context.CancelFunc) {
	go s.waitForShutdown(ctx, cancel)

	log.Info().Msg(fmt.Sprintf("starting server on port %d...", config.Server.Port))

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("server failed to start")
			cancel()
		}
	}()
}

func (s *Server) waitForShutdown(ctx context.Context, cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	log.Info().Msg("shutdown signal received, gracefully stopping server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("server forced to shutdown due to timeout")
	} else {
		log.Info().Msg("server stopped gracefully")
	}

	cancel()
}
