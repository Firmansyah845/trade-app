package app

import (
	"awesomeProjectCr/internal/config"
	"awesomeProjectCr/internal/database"
)

func Init() {
	config.Init()

	database.InitDB()
}

func Shutdown() {
	for _, conn := range database.DBConnection {
		_ = conn.Close()
	}
}
