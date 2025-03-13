package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	// Set global time format
	zerolog.TimeFieldFormat = time.RFC3339

	env := os.Getenv("APP_ENV")
	logLevel := os.Getenv("APP_LOG_LEVEL")

	if env == "dev" {
		// Pretty (human-readable) logs for local dev
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	} else {
		// JSON structured logs for production
		log.Logger = log.Output(os.Stderr)
	}

	level, err := zerolog.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		level = zerolog.InfoLevel // Default to info if parsing fails
	}

	zerolog.SetGlobalLevel(level)
	log.Logger = log.Logger.With().Caller().Logger() // add file name and line number

	log.Info().Str("env", env).Str("level", level.String()).Msg("Logger initialized")
}
