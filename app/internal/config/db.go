package config

type DBConfig struct {
	Host string
	User string
	Pass string
	Name string
	Port string
}

func loadDBConfig() DBConfig {
	return DBConfig{
		Host: getEnv("DB_HOST"),
		User: getEnv("DB_USER"),
		Pass: getEnv("DB_PASS"),
		Name: getEnv("DB_NAME"),
		Port: getEnv("DB_PORT"),
	}
}
