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
	Secret   string
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
			Secret:   getEnv("JWT_SECRET", ""),
		},

		DBConf: DBConfig{
			Host:     getEnv("DB_HOST", ""),
			Port:     getEnv("DB_PORT", ""),
			Username: getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),

			SSLMode:  getEnv("DB_SSLMODE", ""),
			DBName:   getEnv("DB_NAME", ""),
			LogLevel: getEnv("DB_LOGLEVEL", ""),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
