package config

type AppConfig struct {
	Port           string
	Secret         string
	PictureBaseURL string
}

func loadAppConfig() AppConfig {
	return AppConfig{
		Port:           getEnv("APP_PORT"),
		Secret:         getEnv("APP_SECRET"),
		PictureBaseURL: getEnv("APP_PICTURE_BASE_URL"),
	}
}
