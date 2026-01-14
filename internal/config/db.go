package config

import (
	"fmt"
	"time"
)

type DatabaseConfig struct {
	Name                  string
	Host                  string
	User                  string
	Password              string
	Port                  int
	MaxPoolSize           int
	ReadTimeout           time.Duration
	WriteTimeout          time.Duration
	ConnectionMaxLifeTime time.Duration
	SSLMode               string
}

var Database DatabaseConfig

func initDatabaseConfig() {

	Database = DatabaseConfig{
		Name:                  mustGetString("DB_NAME"),
		Host:                  mustGetString("DB_HOST"),
		User:                  mustGetString("DB_USER"),
		Password:              mustGetString("DB_PASSWORD"),
		Port:                  mustGetInt("DB_PORT"),
		MaxPoolSize:           mustGetInt("DB_POOL_SIZE"),
		ReadTimeout:           mustGetDurationMs("DB_READ_TIMEOUT"),
		WriteTimeout:          mustGetDurationMs("DB_WRITE_TIMEOUT"),
		ConnectionMaxLifeTime: mustGetDurationMinute("DB_CONNECTION_MAX_LIFETIME_MINUTE"),
		SSLMode:               mustGetString("SSL_MODE"),
	}
}

func (dc DatabaseConfig) ConnectionURL() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s", dc.User, dc.Password, dc.Host, dc.Port, dc.Name, dc.SSLMode)
}
