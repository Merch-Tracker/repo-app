package helpers

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func ReadBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error("Request error", err)
		return nil, err
	}
	return body, err
}

func GetUUID(w http.ResponseWriter, r *http.Request, pathValue string) (uuid.UUID, error) {
	pathId := r.PathValue(pathValue)
	id, err := uuid.Parse(pathId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.WithField("path value", pathId).Error("Parse user id")
		return uuid.Nil, err
	}
	return id, err
}
