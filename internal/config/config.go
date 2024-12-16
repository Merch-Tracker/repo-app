package config

import "os"

type Config struct {
	HttpConf HttpConfig
	DBConf   DBConfig
}

type HttpConfig struct {
	Host   string
	Port   string
	LogLvl string
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	SSLMode  string
	DBName   string
	LogLvl   string
}

func NewConfig() *Config {
	return &Config{
		HttpConf: HttpConfig{
			Host:   getEnv("HTTP_HOST", "localhost"),
			Port:   getEnv("HTTP_PORT", "9010"),
			LogLvl: getEnv("HTTP_LOGLVL", "Debug"),
		},

		DBConf: DBConfig{
			Host:     getEnv("DB_HOST", "192.168.0.210"),
			Port:     getEnv("DB_PORT", "5432"),
			Username: getEnv("DB_USER", "test_parser"),
			Password: getEnv("DB_PASSWORD", "test_parser"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			DBName:   getEnv("DB_NAME", "parser"),
			LogLvl:   getEnv("DB_LOGLVL", "Info"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
