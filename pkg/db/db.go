package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"repo-app/config"
	"strings"
)

type DB struct {
	*gorm.DB
}

func Connection(c *config.Config) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBConf.Host, c.DBConf.Port, c.DBConf.Username, c.DBConf.Password, c.DBConf.DBName, c.DBConf.SSLMode)

	level := logger.Info

	switch strings.ToLower(c.DBConf.LogLevel) {
	case "silent":
		level = logger.Silent
	case "error":
		level = logger.Error
	case "warn":
		level = logger.Warn
	case "info":
		level = logger.Info
	default:
		level = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: false,
		Logger:                 logger.Default.LogMode(level),
	})

	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
