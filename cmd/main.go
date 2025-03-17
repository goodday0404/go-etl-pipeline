package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/goodday0404/go-etl-pipeline.git/internal/api"
	"github.com/goodday0404/go-etl-pipeline.git/internal/config"
	"github.com/goodday0404/go-etl-pipeline.git/logger"
	"github.com/rs/zerolog/log"
)

func main() {
	logger.InitLogger()
	cfg := loadConfig()
	osSignal := listenToOSsignal()
	server, serverErr := startServer(cfg)
	exitCode := 0

	select {
	case err := <-serverErr:
		log.Error().Caller().Err(err).Msg("❌ Server failed unexpectedly")
		exitCode = 1
	case signal := <-osSignal:
		log.Info().Str("signal", signal.String()).Msg("🛑 Received OS shutdown signal, initiating graceful shutdown...")
	}

	ctx := context.Background()

	if err := server.ShutdownGracefully(ctx); err != nil {
		log.Error().Err(err).Msg("❌ Graceful shutdown failed")
		exitCode = 1
	}

	terminateBackgroundTasks()

	os.Exit(exitCode)
}

func loadConfig() *config.Config {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to load configuration data")
	}

	return cfg
}

func listenToOSsignal() <-chan os.Signal {
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	return stopSignal
}

func startServer(cfg *config.Config) (*api.Server, <-chan error) {
	serverAddress := fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)
	server := api.NewServer(serverAddress, cfg.ShutdownTimeout)
	return server, server.Start()
}

func terminateBackgroundTasks() {
	// TODO: Add background job shutdown logic like ETL pipelines, DB connections, etc.
	log.Info().Msg("✅ All background tasks terminated")
}
