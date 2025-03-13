package main

import (
	"fmt"
	"net/http"

	"github.com/goodday0404/go-etl-pipeline.git/internal/config"
	"github.com/goodday0404/go-etl-pipeline.git/logger"
	"github.com/rs/zerolog/log"
)

func main() {
	logger.InitLogger()

	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration data")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, Go-ETL-Pipeline")
	})

	log.Info().Msgf("Start at http://localhost:%d", cfg.ServerPort)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.DBPort), nil); err != nil {
		log.Fatal().Err(err).Msg("Server stopped")
	}
}
