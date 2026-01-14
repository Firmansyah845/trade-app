package config

import "time"

type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var Server ServerConfig

func initServerConfig() {
	Server = ServerConfig{
		Port:         mustGetInt("SERVER_PORT"),
		ReadTimeout:  mustGetDurationMs("READ_TIMEOUT_MS"),
		WriteTimeout: mustGetDurationMs("WRITE_TIMEOUT_MS"),
	}
}
