package app

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"repo-app/internal/config"
	"repo-app/internal/user"
	"repo-app/pkg/types"
)

type App struct {
	Config     *config.Config
	Router     *http.ServeMux
	HttpServer *Server
	DB         types.Repo
}

func (a *App) Init() {
	a.Router = http.NewServeMux()
	a.HttpServer = NewServer(a.Config, a.Router)

	if a.DB == nil {
		log.Fatal("No database provided")
	}

	// init packages
	NewRootHandler(a.Router)
	user.NewUserHandler(a.Router, a.DB)
}

func (a *App) Start() error {
	log.Info("Starting application")
	err := a.HttpServer.Run()
	if err != nil {
		return err
	}
	return nil
}
