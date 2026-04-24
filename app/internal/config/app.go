package config

type AppConfig struct {
	Port   string
	Secret string
}

func loadAppConfig() AppConfig {
	return AppConfig{
		Port:   getEnv("APP_PORT"),
		Secret: getEnv("APP_SECRET"),
	}
}
