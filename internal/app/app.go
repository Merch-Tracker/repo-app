package app

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"repo-app/internal/auth"
	"repo-app/internal/config"
	"repo-app/internal/images"
	"repo-app/internal/merch"
	"repo-app/internal/user"
	"repo-app/pkg/types"
)

type App struct {
	Config     *config.Config
	Router     *http.ServeMux
	HttpServer *Server
	GrpcServer *grpc.Server
	DB         types.Repo
}

func (a *App) Init() {
	a.Router = http.NewServeMux()
	a.HttpServer = NewHttpServer(a.Config, a.Router)
	a.GrpcServer = NewGrpcServer(a.DB)

	if a.DB == nil {
		log.Fatal("No database provided")
	}

	// init packages
	NewRootHandler(a.Router)
	user.NewUserHandler(a.Router, a.DB)
	auth.NewAuthHandler(a.Router, a.DB)
	merch.NewMerchHandler(a.Router, a.DB)
	images.NewImageHandler(a.Router, a.DB)
}

func (a *App) Start() {
	log.Info("Starting application")

	go func() {
		if err := a.HttpServer.Run(); err != nil {
			log.Fatal(err)
		}
	}()
	log.Debug("HTTP Server started")

	go func() {
		listener, err := net.Listen("tcp", net.JoinHostPort(a.Config.HttpConf.Host, a.Config.HttpConf.GrpcPort))
		if err != nil {
			log.Fatal(err)
		}

		err = a.GrpcServer.Serve(listener)
		if err != nil {
			log.Fatal(err)
		}
	}()
	log.Debug("gRPC Server started")

	select {}
}
