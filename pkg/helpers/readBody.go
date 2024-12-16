package helpers

import (
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
