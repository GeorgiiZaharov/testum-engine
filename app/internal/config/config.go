package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App  AppConfig
	DB   DBConfig
	Ldap LdapConfig
}

func Load() Config {
	_ = godotenv.Load()

	cfg := Config{
		App:  loadAppConfig(),
		DB:   loadDBConfig(),
		Ldap: loadLdapConfig(),
	}

	return cfg
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("missing required env variable: " + key)
	}
	return val
}
