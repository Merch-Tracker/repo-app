package app

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"repo-app/internal/config"
	"repo-app/pkg/middleware"
)

type Server struct {
	Addr   string
	Router *http.ServeMux
}

func NewServer(c *config.Config, router *http.ServeMux) *Server {
	addr := fmt.Sprintf("%s:%s", c.HttpConf.Host, c.HttpConf.Port)
	log.WithFields(log.Fields{"Address": addr}).Debug("Server created")

	return &Server{
		Addr:   addr,
		Router: router,
	}
}

func (s *Server) Run() error {
	mw := []func(http.Handler) http.Handler{
		middleware.CORS,
		middleware.Auth,
	}

	handler := s.mwChain(s.Router, mw)

	server := &http.Server{
		Addr:    s.Addr,
		Handler: handler,
	}

	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) mwChain(handler http.Handler, mw []func(http.Handler) http.Handler) http.Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		handler = mw[i](handler)
	}
	return handler
}
