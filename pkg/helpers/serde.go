package helpers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func SerializeJSON(w http.ResponseWriter, v interface{}) ([]byte, error) {
	marshal, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithField(errMsg, err).Error(serErr)
		return nil, err
	}
	return marshal, nil
}

func DeserializeJSON(w http.ResponseWriter, data []byte, s interface{}) error {
	err := json.Unmarshal(data, &s)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.WithField(errMsg, err).Error(deSerErr)
		return err
	}
	return nil
}
