package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

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
		},
	}

	return server
}
func (s *Server) Start(ctx context.Context, cancel context.CancelFunc) {
	go s.waitForShutDown(ctx, cancel)

	startupMessage := fmt.Sprintf("Starting server on port %s...", strconv.Itoa(config.Server.Port))
	log.Info().Msg(startupMessage)

	go func() {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg(err.Error())
			cancel()
		}
	}()
}

func (s *Server) waitForShutDown(ctx context.Context, cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	log.Info().Msg("stopping server")
	_ = s.server.Shutdown(ctx)

	cancel() // call the cancelFunc to close the shared interrupt channel between REST and gRPC and shutdown both servers
}
