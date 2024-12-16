package app

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"repo-app/internal/config"
)

type Server struct {
	Addr   string
	Router *http.ServeMux
}

func NewServer(c *config.Config) *Server {
	addr := fmt.Sprintf("%s:%s", c.HttpConf.Host, c.HttpConf.Port)
	log.WithFields(log.Fields{"Address": addr}).Debug("Server created")

	return &Server{
		Addr:   addr,
		Router: http.NewServeMux(),
	}
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:    s.Addr,
		Handler: s.Router,
	}

	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
