package user

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"repo-app/pkg/types"
)

type UserHandler struct {
	repo Repo
}

func NewUserHandler(router *http.ServeMux, repo types.Repo) {
	handler := &UserHandler{
		repo: repo,
	}

	err := Migrate(handler.repo)
	if err != nil {
		log.Fatal(migrationErr, err)
	}

	router.HandleFunc("GET /user/", handler.Read())
	router.HandleFunc("PUT /user/", handler.Update())
	router.HandleFunc("DELETE /user/", handler.Delete())
}
