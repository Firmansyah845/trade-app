package config

type AppConfig struct {
	DocsPath string
	AppEnv   string
	Name     string
	Version  string
}

var App AppConfig

func initAppConfig() {
	App.AppEnv = mustGetString("APP_ENV")
	App.Name = mustGetString("APP_NAME")
	App.Version = mustGetString("APP_VERSION")
}
