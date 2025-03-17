package api

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type Server struct {
	shutdownTimeout time.Duration
	*http.Server
}

func NewServer(serverAddress string, shutdownTimeout time.Duration) *Server {
	return &Server{
		shutdownTimeout: shutdownTimeout,

		Server: &http.Server{
			Addr:         serverAddress,
			Handler:      newRouter(),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  30 * time.Second,
		}}
}

func (s Server) Start() <-chan error {
	serverErr := make(chan error, 1)

	go func() {
		log.Info().Msgf("🚀 Server listening on %s", s.Addr)

		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Caller().Err(err).Msg("❌ Failed to start server")
			serverErr <- err
		}
	}()

	return serverErr
}

func (s Server) ShutdownGracefully(ctx context.Context) error {
	log.Info().Msgf("🛑 Starting graceful shutdown. Waiting up to %s to finish current requests...", s.shutdownTimeout)

	timeoutCtx, cancel := context.WithTimeout(ctx, s.shutdownTimeout)
	defer cancel()

	if err := s.Shutdown(timeoutCtx); err != nil {
		log.Error().Caller().Err(err).Msg("❌ Failed to shutdown the server gracefully")
		return err
	}

	log.Info().Msg("✅ Sever is successfully closed")
	return nil
}
