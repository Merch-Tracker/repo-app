package merch

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"repo-app/pkg/helpers"
)

type Label struct {
	gorm.Model
	LabelUuid uuid.UUID
	UserUuid  uuid.UUID
	Name      string `json:"name"`
	Color     string `json:"color"`
	BgColor   string `json:"bg_color"`
}

type CardLabel struct {
	LabelUuid uuid.UUID `json:"label_uuid"`
	UserUuid  uuid.UUID
	MerchUuid uuid.UUID `json:"merch_uuid"`
}

func (m *MerchHandler) NewLabel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := helpers.ReadBody(w, r)
		if err != nil {
			return
		}

		newLabel := &Label{}
		err = helpers.DeserializeJSON(w, body, newLabel)
		if err != nil {
			return
		}

		if newLabel.Color == "" {
			newLabel.Color = "#000000"
		}

		if newLabel.BgColor == "" {
			newLabel.BgColor = "#ffffff"
		}

		newLabel.UserUuid = helpers.GetUserUuid(r)
		newLabel.LabelUuid = uuid.New()

		err = newLabel.Create(m.repo)
		if err != nil {
			http.Error(w, labelsCreateError, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(labelsCreateError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		log.Info(labelsCreateSuccess)
	}
}

func (m *MerchHandler) GetLabels() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := Label{UserUuid: helpers.GetUserUuid(r)}

		labels, err := l.ReadAll(m.repo)
		if err != nil {
			http.Error(w, labelsGetAllError, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(labelsGetAllError)
			return
		}

		response, err := helpers.SerializeJSON(w, labels)
		if err != nil {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		log.Info(labelsGetSuccess)
	}
}

func (m *MerchHandler) UpdateLabel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		labelID, err := helpers.GetPathUuid(w, r, "label_uuid")
		if err != nil {
			log.Error(labelsGetIdError)
			return
		}

		body, err := helpers.ReadBody(w, r)
		if err != nil {
			return
		}

		updateLabel := &Label{}
		err = helpers.DeserializeJSON(w, body, updateLabel)
		if err != nil {
			return
		}

		updateLabel.UserUuid = helpers.GetUserUuid(r)
		updateLabel.LabelUuid = labelID

		err = updateLabel.Update(m.repo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(labelsUpdateError)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info(labelsUpdateSuccess)
	}
}

func (m *MerchHandler) DeleteLabel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		labelID, err := helpers.GetPathUuid(w, r, "label_uuid")
		if err != nil {
			log.Error(labelsGetIdError)
			return
		}

		deleteLabel := &Label{
			UserUuid:  helpers.GetUserUuid(r),
			LabelUuid: labelID,
		}

		err = deleteLabel.Delete(m.repo)
		if err != nil {
			http.Error(w, labelsDeleteError, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(labelsDeleteError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		log.Info(labelsDeleteSuccess)
	}
}

func (m *MerchHandler) AttachLabel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := helpers.ReadBody(w, r)
		if err != nil {
			log.Error(labelAttachError)
			return
		}

		cardLabel := &CardLabel{}
		err = helpers.DeserializeJSON(w, body, cardLabel)
		if err != nil {
			log.Error(labelAttachError)
			return
		}

		cardLabel.UserUuid = helpers.GetUserUuid(r)
		err = cardLabel.Attach(m.repo)
		if err != nil {
			http.Error(w, labelAttachError, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(labelAttachError)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info(labelAttachSuccess)
	}
}

func (m *MerchHandler) DetachLabel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := helpers.ReadBody(w, r)
		if err != nil {
			log.Error(labelDetachError)
			return
		}

		cardLabel := &CardLabel{}
		err = helpers.DeserializeJSON(w, body, cardLabel)
		if err != nil {
			log.Error(labelDetachError)
			return
		}

		cardLabel.UserUuid = helpers.GetUserUuid(r)
		err = cardLabel.Detach(m.repo)
		if err != nil {
			http.Error(w, labelDetachError, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(labelDetachError)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info(labelDetachSuccess)
	}
}

// Inner methods
func (l *Label) Create(repo Repo) error {
	err := repo.Create(l)
	if err != nil {
		return err
	}
	return nil
}

func (l *Label) ReadAll(repo Repo) (*[]Label, error) {
	labels := &[]Label{}
	params := make(map[string]any)
	params["user_uuid"] = l.UserUuid

	err := repo.ReadManySimple(labels, params)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

func (l *Label) Update(repo Repo) error {
	params := make(map[string]any)
	params["user_uuid"] = l.UserUuid
	params["label_uuid"] = l.LabelUuid

	err := repo.Update(l, params)
	if err != nil {
		return err
	}
	return nil
}

func (l *Label) Delete(repo Repo) error {
	params := make(map[string]any)
	params["user_uuid"] = l.UserUuid
	params["label_uuid"] = l.LabelUuid

	err := repo.Delete(l, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *CardLabel) Attach(repo Repo) error {
	err := repo.Create(c)
	if err != nil {
		return err
	}
	return nil
}

func (c *CardLabel) Detach(repo Repo) error {
	params := make(map[string]any)
	params["user_uuid"] = c.UserUuid
	params["merch_uuid"] = c.MerchUuid
	params["label_uuid"] = c.LabelUuid

	err := repo.Delete(c, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *CardLabel) ReadAll(repo Repo) (*[]CardLabel, error) {
	params := make(map[string]any)
	params["user_uuid"] = c.UserUuid

	labels := &[]CardLabel{}
	err := repo.ReadManySimple(labels, params)
	if err != nil {
		return nil, err
	}

	return labels, nil
}
