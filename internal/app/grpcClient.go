package app

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"repo-app/config"
	pb "repo-app/pkg/pricewatcher"
)

func NewGrpcClient(c *config.Config) pb.PriceWatcherClient {
	var opts []grpc.DialOption
	insec := grpc.WithTransportCredentials(insecure.NewCredentials())
	opts = append(opts, insec)

	conn, err := grpc.NewClient(c.HttpConf.Host+":"+c.HttpConf.GrpcClientPort, opts...)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("gRPC Client created")

	return pb.NewPriceWatcherClient(conn)
}
