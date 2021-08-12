package geaconfig

type config struct {
	Version  string
	SiteName string
}

var configInstance *config

func GEAConfig() *config {
	if configInstance != nil {
		return configInstance
	}
	configInstance = new(config)
	return configInstance
}
