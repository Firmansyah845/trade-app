package database

import (
	"awesomeProjectCr/internal/config"
	"awesomeProjectCr/pkg/db"
	"database/sql"

	"github.com/rs/zerolog/log"
)

const (
	PostgresDb = "postgresDB"
)

var DBConnection = map[string]*sql.DB{}

func InitDB() {
	dbConfig := db.Config{
		Driver:          "postgres",
		URL:             config.Database.ConnectionURL(),
		MaxIdleConns:    config.Database.MaxPoolSize,
		MaxOpenConns:    config.Database.MaxPoolSize,
		ConnMaxLifeTime: config.Database.ConnectionMaxLifeTime,
	}

	postgresDB, err := db.NewDB(&dbConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres database")
	}

	DBConnection[PostgresDb] = postgresDB

	log.Info().
		Str("database", PostgresDb).
		Int("max_pool_size", config.Database.MaxPoolSize).
		Msg("database connection established")
}
