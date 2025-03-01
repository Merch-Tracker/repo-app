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
	"repo-app/internal/notify"
	"repo-app/internal/status"
	"repo-app/internal/user"
	"repo-app/pkg/jwt"
	pb "repo-app/pkg/pricewatcher"
	"repo-app/pkg/types"
)

type App struct {
	Config          *config.Config
	DB              types.Repo
	Router          *http.ServeMux
	HttpServer      *Server
	GrpcServer      *grpc.Server
	GrpcClient      pb.PriceWatcherClient
	NotifierService *notify.NotificationService
	NotifierChan    chan struct{}
}

func (a *App) Init() {
	// init vars
	a.NotifierChan = make(chan struct{}, 30)
	jwt.Secret = a.Config.HttpConf.Secret

	if a.DB == nil {
		log.Fatal(noDBErr)
	}

	// init services
	a.Router = http.NewServeMux()
	a.HttpServer = NewHttpServer(a.Config, a.Router)
	a.GrpcServer = NewGrpcServer(a.DB, a.NotifierChan)
	a.GrpcClient = NewGrpcClient(a.Config)
	a.NotifierService = notify.NewNotificationService(a.DB, a.NotifierChan)

	// init packages
	NewRootHandler(a.Router)
	user.NewUserHandler(a.Router, a.DB)
	auth.NewAuthHandler(a.Router, a.DB)
	merch.NewMerchHandler(a.Router, a.DB)
	images.NewImageHandler(a.Router, a.DB)
	notify.NewNotifierHandler(a.DB, a.Router)
	status.NewStatusHandler(a.Router, a.GrpcClient)
}

func (a *App) Start() {
	log.Info(appStart)

	go func() {
		if err := a.HttpServer.Run(); err != nil {
			log.WithField(errMsg, err).Fatal(httpServerFatal)
		}
	}()
	log.Info(httpServerStart)

	go func() {
		listener, err := net.Listen("tcp", net.JoinHostPort(a.Config.HttpConf.Host, a.Config.HttpConf.GrpcServerPort))
		if err != nil {
			log.WithField(errMsg, err).Fatal(grpcServerFatal)
		}

		err = a.GrpcServer.Serve(listener)
		if err != nil {
			log.WithField(errMsg, err).Fatal(grpcServerFatal)
		}
	}()
	log.Info(grpcServerStart)

	go func() {
		if err := a.NotifierService.Run(); err != nil {
			log.WithField(errMsg, err).Error(notificationServiceError)
		}
	}()
	log.Info(notificationServiceSuccess)

	select {}
}
