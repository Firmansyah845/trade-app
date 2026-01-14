package config

import "os"

type ApmConfig struct {
	Active      string
	ApiKey      string
	Env         string
	Url         string
	ServiceName string
}

var ApmDataConfig ApmConfig

func initApmConfig() {

	ApmDataConfig.Active = mustGetString("ELASTIC_APM_ACTIVE")
	ApmDataConfig.ApiKey = mustGetString("ELASTIC_APM_API_KEY")
	ApmDataConfig.Env = mustGetString("ELASTIC_APM_ENVIRONMENT")
	ApmDataConfig.Url = mustGetString("ELASTIC_APM_SERVER_URL")
	ApmDataConfig.ServiceName = mustGetString("ELASTIC_APM_SERVICE_NAME")

	_ = os.Setenv("ELASTIC_APM_ACTIVE", ApmDataConfig.Active)
	_ = os.Setenv("ELASTIC_APM_API_KEY", ApmDataConfig.ApiKey)
	_ = os.Setenv("ELASTIC_APM_ENVIRONMENT", ApmDataConfig.Env)
	_ = os.Setenv("ELASTIC_APM_SERVER_URL", ApmDataConfig.Url)
	_ = os.Setenv("ELASTIC_APM_SERVICE_NAME", ApmDataConfig.ServiceName)

}
