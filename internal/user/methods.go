package user

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"repo-app/pkg/helpers"
	"repo-app/pkg/types"
)

func (u *UserHandler) Read() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		readUser := &User{UserUUID: helpers.GetUUID(r)}
		payload := &Personal{}

		err := readUser.ReadOnePayload(u.repo, payload)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.WithField("msg", err).Error(types.UserReadError)
			return
		}

		response, err := helpers.SerializeJSON(&w, payload)
		if err != nil {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		log.WithField("response", payload).Info(types.UserReadSuccess)
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
