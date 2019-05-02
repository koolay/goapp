package gateway

type Configuration struct {
	LogLevel   string
	ConfigFile string
	Plugin     struct {
		Dir     string
		Pattern string
	}
}

func DefaultConfiguration() Configuration {
	return Configuration{
		LogLevel:   "ERROR",
		ConfigFile: "gateway.json",
		Plugin: struct {
			Dir     string
			Pattern string
		}{
			Dir:     "./plugins/",
			Pattern: ".so",
		},
	}
}
