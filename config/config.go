package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port               string
	MaxConcurrentSends int
	SendTimeout        time.Duration
	EnvLogLevel        string
}

func Load() *Config {
	// defaults
	cfg := &Config{
		Port:               getEnv("PORT", "3000"),
		MaxConcurrentSends: getEnvInt("MAX_CONCURRENT_SENDS", 10),
		SendTimeout:        time.Duration(getEnvInt("SEND_TIMEOUT_SECONDS", 30)) * time.Second,
		EnvLogLevel:        getEnv("LOG_LEVEL", "info"),
	}
	return cfg
}

func getEnv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func getEnvInt(k string, d int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return d
}
