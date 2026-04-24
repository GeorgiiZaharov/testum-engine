package config

import (
	"log"
	"strconv"
	"time"
)

type LdapConfig struct {
	Server  string
	Port    int
	BaseDN  string
	Timeout time.Duration
}

func loadLdapConfig() LdapConfig {
	port, _ := strconv.Atoi(getEnv("LDAP_PORT"))

	timeout, err := time.ParseDuration(getEnv("LDAP_TIMEOUT"))
	if err != nil {
		log.Fatal("Wrong duration")
		timeout = 5 * time.Second // дефолт
	}

	return LdapConfig{
		Server:  getEnv("LDAP_SERVER"),
		Port:    port,
		BaseDN:  getEnv("LDAP_BASE_DN"),
		Timeout: timeout,
	}
}
