package helpers

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"repo-app/pkg/jwt"
	"repo-app/pkg/types"
)

func ReadBody(w *http.ResponseWriter, r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		(*w).WriteHeader(http.StatusBadRequest)
		log.Error("Request error", err)
		return nil, err
	}
	return body, nil
}

func GetUserUuid(r *http.Request) uuid.UUID {
	return r.Context().Value(types.UserDataKey).(*jwt.Data).UserID
}

func GetPathUuid(w *http.ResponseWriter, r *http.Request, pathValue string) (uuid.UUID, error) {
	pathUuid, err := uuid.Parse(r.PathValue(pathValue))
	if err != nil {
		(*w).WriteHeader(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"error": err,
			"uuid":  pathValue,
		}).Error("Request error")
		return uuid.Nil, err
	}
	return pathUuid, nil
}
