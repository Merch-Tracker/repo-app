package app

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"repo-app/internal/merch"
	pb "repo-app/pkg/pricewatcher"
	"repo-app/pkg/types"
	"time"
)

type Repo types.Repo

type pwServer struct {
	pb.UnimplementedPriceWatcherServer
	repo Repo
}

func NewGrpcServer(repo Repo) *grpc.Server {
	grpcServer := grpc.NewServer()
	pwServ := pwServer{repo: repo}
	pb.RegisterPriceWatcherServer(grpcServer, &pwServ)
	return grpcServer
}

func (s *pwServer) GetMerch(req *emptypb.Empty, stream pb.PriceWatcher_GetMerchServer) error {
	me := merch.GMerch{}
	merchList, err := me.ReadAll(s.repo)
	if err != nil {
		log.WithField(errMsg, err).Error(grpcGetMerchRepoReadError)
		return err
	}

	for _, m := range merchList {
		response := &pb.MerchRequest{
			MerchUuid:    m.MerchUuid.String(),
			Link:         m.Link,
			ParseTag:     m.ParseTag,
			ParseSubs:    m.ParseSubstring,
			CookieValues: m.CookieValues,
			Separator:    m.Separator,
		}

		err = stream.Send(response)
		if err != nil {
			log.WithField(errMsg, err).Error(grpcGetMerchStreamError)
			return err
		}
		log.WithField(respMsg, response).Debug(grpcGetMerchSuccess)
	}
	return nil
}

func (s *pwServer) PostMerch(stream pb.PriceWatcher_PostMerchServer) error {
	saveInterval := time.Second * 2
	batch := make([]merch.MerchInfo, 0)

	ticker := time.NewTicker(saveInterval)
	defer ticker.Stop()

	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if len(batch) > 0 {
					err := s.SaveToDB(batch)
					if err != nil {
						log.WithField(errMsg, err).Error(grpcPostMerchBatchError)
					}
				}
			}
		}
	}()

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			log.Debug(grpcEOF)
			break
		}

		if err != nil {
			log.WithField(errMsg, err).Error(grpcReceiveError)
			return err
		}

		entry := merch.MerchInfo{MerchUuid: uuid.MustParse(response.MerchUuid), Price: uint(response.Price)}
		batch = append(batch, entry)
		log.WithField(respMsg, entry).Debug(grpcReceiveSuccess)
	}

	close(done)
	if len(batch) > 0 {
		err := s.SaveToDB(batch)
		if err != nil {
			log.WithField(errMsg, err).Error(grpcPostMerchBatchError)
			return err
		}
	}

	return nil
}

func (s *pwServer) SaveToDB(list []merch.MerchInfo) error {
	err := s.repo.Save(&list)
	if err != nil {
		return err
	}
	return nil
}
