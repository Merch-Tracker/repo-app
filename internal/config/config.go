package config

import "os"

type Config struct {
	HttpConf ServerConfig
	DBConf   DBConfig
}

type ServerConfig struct {
	Host     string
	HttpPort string
	GrpcPort string
	LogLvl   string
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	SSLMode  string
	DBName   string
	LogLevel string
}

func NewConfig() *Config {
	return &Config{
		HttpConf: ServerConfig{
			Host:     getEnv("HTTP_HOST", "localhost"),
			HttpPort: getEnv("HTTP_PORT", "9010"),
			GrpcPort: getEnv("GRPC_PORT", "9050"),
			LogLvl:   getEnv("HTTP_LOGLEVEL", "Debug"),
		},

		DBConf: DBConfig{
			Host:     getEnv("DB_HOST", "192.168.0.210"),
			Port:     getEnv("DB_PORT", "5432"),
			Username: getEnv("DB_USER", "test_parser"),
			Password: getEnv("DB_PASSWORD", "test_parser"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			DBName:   getEnv("DB_NAME", "parser"),
			LogLevel: getEnv("DB_LOGLEVEL", "Info"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
