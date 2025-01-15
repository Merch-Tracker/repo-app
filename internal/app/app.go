package app

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"repo-app/config"
	"repo-app/internal/auth"
	"repo-app/internal/images"
	"repo-app/internal/merch"
	"repo-app/internal/user"
	"repo-app/pkg/jwt"
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
		log.Fatal(noDBErr)
	}

	// init vars
	jwt.Secret = a.Config.HttpConf.Secret

	// init packages
	NewRootHandler(a.Router)
	user.NewUserHandler(a.Router, a.DB)
	auth.NewAuthHandler(a.Router, a.DB)
	merch.NewMerchHandler(a.Router, a.DB)
	images.NewImageHandler(a.Router, a.DB)
}

func (a *App) Start() {
	log.Info(appStart)

	go func() {
		if err := a.HttpServer.Run(); err != nil {
			log.WithField(errMsg, err).Fatal(httpServerFatal)
		}
	}()
	log.Debug(httpServerStart)

	go func() {
		listener, err := net.Listen("tcp", net.JoinHostPort(a.Config.HttpConf.Host, a.Config.HttpConf.GrpcPort))
		if err != nil {
			log.WithField(errMsg, err).Fatal(grpcServerFatal)
		}

		err = a.GrpcServer.Serve(listener)
		if err != nil {
			log.WithField(errMsg, err).Fatal(grpcServerFatal)
		}
	}()
	log.Debug(grpcServerStart)

	select {}
}
