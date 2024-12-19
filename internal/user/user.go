package user

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"repo-app/pkg/helpers"
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
		log.Fatal("Migration error", err)
	}

	router.HandleFunc("GET /user/", handler.Read())
	router.HandleFunc("PUT /user/{userid}", handler.Update())
	router.HandleFunc("DELETE /user/{userid}", handler.Delete())
}

func (u *UserHandler) Read() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		readUser := &User{}
		readUser.UserUUID = helpers.GetUUID(r)
		err := readUser.ReadOne(u.repo)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.WithField("msg", err).Error("Read user")
			return
		}

		response, err := helpers.SerializeJSON(&w, readUser)
		if err != nil {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func (u *UserHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := helpers.ReadBody(&w, r)
		if err != nil {
			return
		}

		updateUser := &User{}
		err = helpers.DeserializeJSON(&w, body, &updateUser)
		if err != nil {
			return
		}

		updateUser.UserUUID = helpers.GetUUID(r)
		err = updateUser.Update(u.repo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Update user")
			return
		}

		w.WriteHeader(http.StatusNoContent)
		log.WithFields(log.Fields{"data": updateUser}).Info("Update user")
	}
}

func (u *UserHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deleteUser := &User{UserUUID: helpers.GetUUID(r)}
		err := deleteUser.Delete(u.repo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Delete user")
			return
		}
		w.WriteHeader(http.StatusNoContent)
		log.WithFields(log.Fields{"data": deleteUser}).Info("Delete user")
	}
}
