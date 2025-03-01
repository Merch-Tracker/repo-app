package config

import "os"

type Config struct {
	HttpConf ServerConfig
	DBConf   DBConfig
}

type ServerConfig struct {
	Host           string
	HttpPort       string
	GrpcServerPort string
	GrpcClientPort string
	LogLvl         string
	Secret         string
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
			Host:           getEnv("HTTP_HOST", "localhost"),
			HttpPort:       getEnv("HTTP_PORT", "9010"),
			GrpcServerPort: getEnv("GRPC_SERVER_PORT", "9050"),
			GrpcClientPort: getEnv("GRPC_CLIENT_PORT", "9060"),
			LogLvl:         getEnv("HTTP_LOGLEVEL", "Info"),
			Secret:         getEnv("JWT_SECRET", ""),
		},

		DBConf: DBConfig{
			Host:     getEnv("DB_HOST", ""),
			Port:     getEnv("DB_PORT", ""),
			Username: getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),

			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			DBName:   getEnv("DB_NAME", ""),
			LogLevel: getEnv("DB_LOGLEVEL", "silent"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
