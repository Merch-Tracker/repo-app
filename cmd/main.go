package main

import (
	log "github.com/sirupsen/logrus"
	"repo-app/internal/app"
	"repo-app/internal/config"
	"repo-app/internal/logging"
	"repo-app/pkg/db"
)

func main() {
	c := config.NewConfig()
	logging.LogSetup(c.HttpConf.LogLvl)

	database, err := db.Connection(c)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Database connection failed")
	}

	appl := app.App{
		Config: c,
		DB:     database,
	}

	appl.Init()

	err = appl.Start()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Application start error")
	}
}
