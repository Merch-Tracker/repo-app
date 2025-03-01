package status

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"repo-app/pkg/helpers"
	pb "repo-app/pkg/pricewatcher"
	"time"
)

type StatusHandler struct {
	client pb.PriceWatcherClient
}

type parserStatus struct {
	Origin      string    `json:"origin"`
	StartTime   time.Time `json:"start_time"`
	CheckPeriod uint32    `json:"check_period"`
	LastCheck   time.Time `json:"last_check"`
	NextCheck   time.Time `json:"next_check"`
	NumCpus     uint32    `json:"num_cpus"`
}

func NewStatusHandler(router *http.ServeMux, rc pb.PriceWatcherClient) {
	handler := &StatusHandler{
		client: rc,
	}

	router.HandleFunc("GET /status/parser", handler.GetParserStatus())
}

func (s *StatusHandler) GetParserStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Context(context.Background())
		req := pb.StatusRequest{}

		resp, err := s.client.ParserInfo(ctx, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error("Error getting parser info by gRPC")
			return
		}

		lastCheck := time.Unix(int64(resp.LastCheck), 0)

		rsp := parserStatus{
			Origin:      "surugaya",
			StartTime:   time.Unix(int64(resp.StartTime), 0),
			CheckPeriod: resp.CheckPeriod,
			LastCheck:   lastCheck,
			NextCheck:   lastCheck.Add(time.Hour * time.Duration(resp.CheckPeriod)),
			NumCpus:     resp.NumCpus,
		}

		response, err := helpers.SerializeJSON(w, &rsp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(response)
		log.Info("Getting parser status success")
	}
}
