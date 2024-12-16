package app

import (
	log "github.com/sirupsen/logrus"
	"repo-app/pkg/types"
)

type App struct {
	HttpServer *Server
	DBW        *types.Repo
}

func (a *App) Start() error {
	log.Info("Starting application")
	err := a.HttpServer.Run()
	if err != nil {
		return err
	}
	return nil
}
