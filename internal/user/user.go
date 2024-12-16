package user

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"repo-app/pkg/helpers"
	"repo-app/pkg/types"
	"strings"
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
		log.Fatal("Migration error", err)
	}

	router.HandleFunc("POST /user", handler.Create())
	router.HandleFunc("GET /user/{userid}", handler.Read())
	router.HandleFunc("PUT /user/{userid}", handler.Update())
	router.HandleFunc("DELETE /user/{userid}", handler.Delete())
}

func (u *UserHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		newUser := User{}

		body, err := helpers.ReadBody(w, r)
		if err != nil {
			return
		}
		log.WithFields(log.Fields{"data": string(body)}).Debug("Request body")

		err = helpers.DeserializeJSON(w, body, &newUser)
		if err != nil {
			return
		}
		log.WithFields(log.Fields{"data": newUser}).Debug("Deserialized")

		err = newUser.Create(u.repo)
		switch {
		case err == nil:
			w.WriteHeader(http.StatusCreated)
			return

		case strings.Contains(err.Error(), "duplicate key value violates unique constraint"):
			w.WriteHeader(http.StatusConflict)
			log.Error("User already exists")
			return

		default:
			w.WriteHeader(http.StatusInternalServerError)
			log.WithField("msg", err).Error("Create user")
			return
		}
	}
}

func (u *UserHandler) Read() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (u *UserHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (u *UserHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
