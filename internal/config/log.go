package config

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func initLogConfig() {
	levelStr := strings.ToLower(mustGetString("LOG_LEVEL"))

	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs // lebih presisi, lebih ringan dari ISO8601

	env := strings.ToLower(mustGetString("APP_ENV"))

	var logger zerolog.Logger
	if env == "development" || env == "local" {
		// pretty print untuk development
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
			With().
			Timestamp().
			Caller().
			Str("service", mustGetString("APP_NAME")).
			Logger()
	} else {
		// pure JSON untuk production/staging
		logger = zerolog.New(os.Stdout).
			With().
			Timestamp().
			Caller().
			Str("service", mustGetString("APP_NAME")).
			Str("env", env).
			Logger()
	}

	log.Logger = logger

	log.Info().
		Str("level", level.String()).
		Str("env", env).
		Msg("logger initialized")
}
