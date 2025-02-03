package helpers

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"repo-app/pkg/jwt"
)

func ReadBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, requestError, http.StatusBadRequest)
		log.WithField(errMsg, err).Error(requestError)
		return nil, err
	}
	return body, nil
}

func GetUserUuid(r *http.Request) uuid.UUID {
	return r.Context().Value(jwt.UserDataKey).(*jwt.Data).UserID
}

func GetPathUuid(w http.ResponseWriter, r *http.Request, pathValue string) (uuid.UUID, error) {
	pathUuid, err := uuid.Parse(r.PathValue(pathValue))
	if err != nil {
		http.Error(w, requestError, http.StatusBadRequest)
		log.WithFields(log.Fields{
			errMsg: err,
			idMsg:  pathValue,
		}).Error(requestError)
		return uuid.Nil, err
	}
	return pathUuid, nil
}
