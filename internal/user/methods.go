package user

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"repo-app/pkg/helpers"
)

func (u *UserHandler) Read() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		readUser := &User{UserUUID: helpers.GetUserUuid(r)}
		payload := &Personal{}

		err := readUser.ReadOnePayload(u.repo, payload)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.WithField(errMsg, err).Error(userReadError)
			return
		}

		response, err := helpers.SerializeJSON(&w, payload)
		if err != nil {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		log.WithField(respMsg, payload).Info(userReadSuccess)
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

		updateUser.UserUUID = helpers.GetUserUuid(r)
		err = updateUser.Update(u.repo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(userUpdateError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		log.WithField(noRespMsg, updateUser).Info(userUpdateSuccess)
	}
}

func (u *UserHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deleteUser := &User{UserUUID: helpers.GetUserUuid(r)}
		err := deleteUser.Delete(u.repo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(userDeleteError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		log.WithField(noRespMsg, deleteUser).Info(userDeleteSuccess)
	}
}
