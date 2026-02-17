package app

import (
	"awesomeProjectCr/internal/config"
	"awesomeProjectCr/internal/database"

	"github.com/rs/zerolog/log"
)

func Init() {
	config.Init()
	database.InitDB()
}

func Shutdown() {
	for name, conn := range database.DBConnection {
		if err := conn.Close(); err != nil {
			log.Error().
				Err(err).
				Str("database", name).
				Msg("failed to close database connection")
			continue
		}
		log.Info().
			Str("database", name).
			Msg("database connection closed")
	}
}
