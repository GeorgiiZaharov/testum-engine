package config

type DBConfig struct {
	Path string
}

func loadDBConfig() DBConfig {
	return DBConfig{
		Path: getEnv("DB_PATH"),
	}
}
