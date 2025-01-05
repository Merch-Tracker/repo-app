package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"repo-app/config"
)

type DB struct {
	*gorm.DB
}

func Connection(c *config.Config) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBConf.Host, c.DBConf.Port, c.DBConf.Username, c.DBConf.Password, c.DBConf.DBName, c.DBConf.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: false,
		Logger:                 logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
