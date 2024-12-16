package main

import (
	log "github.com/sirupsen/logrus"
	"repo-app/internal/app"
	"repo-app/internal/config"
	"repo-app/internal/logging"
	"repo-app/pkg/db"
	"repo-app/pkg/types"
)

func main() {
	c := config.NewConfig()
	logging.LogSetup(c.HttpConf.LogLvl)

	database, err := db.NewDB(c)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Database connection failed")
	}

	var dbw types.Repo = database

	appl := app.App{
		HttpServer: app.NewServer(c),
		DBW:        &dbw,
	}

	err = appl.Start()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Application start error")
	}
}
