package notify

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"repo-app/pkg/helpers"
	"repo-app/pkg/types"
)

type Repo types.Repo

type NotifierHandler struct {
	repo   Repo
	router *http.ServeMux
}

func NewNotifierHandler(repo Repo, router *http.ServeMux) {
	var err error
	handler := &NotifierHandler{
		repo:   repo,
		router: router,
	}

	err = MigrateNotifiers(handler.repo)
	if err != nil {
		log.WithField(errMsg, err).Fatal(migrationNotifiersError)
	}

	err = MigrateNotifyMessages(handler.repo)
	if err != nil {
		log.WithField(errMsg, err).Fatal(migrationNotifyMessagesError)
	}

	log.Debug(migrationSuccess)

	router.HandleFunc("GET /notifications", handler.GetNotifications())
	router.HandleFunc("POST /notifications", handler.MarkAsRead())
}

func (n *NotifierHandler) GetNotifications() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usr := helpers.GetUserUuid(r)

		nm := NotifyMessageResponse{UserUuid: usr.String()}
		messages, err := nm.ReadAll(n.repo)
		if err != nil {
			http.Error(w, readNotificationError, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(readNotificationError)
			return
		}

		if len(*messages) == 0 {
			w.WriteHeader(http.StatusNoContent)
			log.WithField(respMsg, http.StatusNoContent).Debug(noContent)
			return
		}

		response, err := helpers.SerializeJSON(&w, messages)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
		log.WithField(respMsg, len(response)).Info(readNotificationSuccess)
	}
}

func (n *NotifierHandler) MarkAsRead() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usr := helpers.GetUserUuid(r)

		body, err := helpers.ReadBody(&w, r)
		if err != nil {
			return
		}

		markList := &[]NotifyMessage{}

		err = helpers.DeserializeJSON(&w, body, markList)
		if err != nil {
			return
		}

		for _, message := range *markList {
			message.UserUuid = usr.String()
			message.Viewed = true
		}

		nm := NotifyMessage{}
		err = nm.MarkAsRead(n.repo, markList)
		if err != nil {
			http.Error(w, markAsReadError, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(markAsReadError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		log.WithField(respMsg, http.StatusNoContent).Info(markAsReadSuccess)
	}
}
