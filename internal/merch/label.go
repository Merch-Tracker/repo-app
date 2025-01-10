package merch

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"repo-app/pkg/helpers"
)

type Label struct {
	gorm.Model
	LabelUuid uuid.UUID
	OwnerUuid uuid.UUID
	Name      string `json:"name"`
	Color     string `json:"color"`
	BgColor   string `json:"bg_color"`
}

type CardLabel struct {
	LabelUuid uuid.UUID `json:"label_uuid"`
	OwnerUuid uuid.UUID
	MerchUuid uuid.UUID `json:"merch_uuid"`
}

func MigrateLabel(repo Repo) error {
	err := repo.Migrate(Label{})
	if err != nil {
		return err
	}
	return nil
}

func MigrateCardLabel(repo Repo) error {
	err := repo.Migrate(CardLabel{})
	if err != nil {
		return err
	}
	return nil
}

// HTTP Handlers
func (m *MerchHandler) NewLabel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := helpers.ReadBody(&w, r)
		if err != nil {
			log.Error("Error in new label method")
			return
		}

		newLabel := &Label{}
		err = helpers.DeserializeJSON(&w, body, newLabel)
		if err != nil {
			log.Error("Error in deserialize label method")
			return
		}

		if newLabel.Color == "" {
			newLabel.Color = "#000000"
		}

		if newLabel.BgColor == "" {
			newLabel.BgColor = "#ffffff"
		}

		newLabel.OwnerUuid = helpers.GetUserUuid(r)
		newLabel.LabelUuid = uuid.New()

		err = newLabel.Create(m.repo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error("Error creating new label")
			return
		}

		w.WriteHeader(http.StatusCreated)
		log.Info("New label created")
	}
}

func (m *MerchHandler) GetLabel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := Label{OwnerUuid: helpers.GetUserUuid(r)}

		labels, err := l.ReadAll(m.repo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error("Error in get label")
			return
		}

		response, err := helpers.SerializeJSON(&w, labels)
		if err != nil {
			log.Error("Error in get label")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		log.Info("Get label success")
	}
}

func (m *MerchHandler) UpdateLabel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		labelID, err := helpers.GetPathUuid(&w, r, "label_uuid")
		if err != nil {
			log.Error("Error in get label uuid")
			return
		}

		body, err := helpers.ReadBody(&w, r)
		if err != nil {
			log.Error("Error in update label method")
			return
		}

		fmt.Println("BODY", string(body))

		updateLabel := &Label{}
		err = helpers.DeserializeJSON(&w, body, updateLabel)
		if err != nil {
			log.Error("Error in deserialize label method")
			return
		}

		updateLabel.OwnerUuid = helpers.GetUserUuid(r)
		updateLabel.LabelUuid = labelID
		err = updateLabel.Update(m.repo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error("Error in update label")
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info("Update label success")
	}
}
func (m *MerchHandler) DeleteLabel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		labelID, err := helpers.GetPathUuid(&w, r, "label_uuid")
		if err != nil {
			log.Error("Error in get delete label uuid")
			return
		}

		deleteLabel := &Label{
			OwnerUuid: helpers.GetUserUuid(r),
			LabelUuid: labelID,
		}

		err = deleteLabel.Delete(m.repo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error("Error in delete deleteLabel")
			return
		}
		w.WriteHeader(http.StatusNoContent)
		log.Info("Delete delete label success")
	}
}

func (m *MerchHandler) AttachLabel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := helpers.ReadBody(&w, r)
		if err != nil {
			log.Error("Error in attach label method")
			return
		}

		cardLabel := &CardLabel{}
		err = helpers.DeserializeJSON(&w, body, cardLabel)
		if err != nil {
			log.Error("Error in deserialize label method")
			return
		}

		cardLabel.OwnerUuid = helpers.GetUserUuid(r)
		err = cardLabel.Attach(m.repo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error("Error in attach label")
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info("Attach label success")
	}
}

func (m *MerchHandler) DetachLabel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := helpers.ReadBody(&w, r)
		if err != nil {
			log.Error("Error in detach label method")
			return
		}

		cardLabel := &CardLabel{}
		err = helpers.DeserializeJSON(&w, body, cardLabel)
		if err != nil {
			log.Error("Error in deserialize label method")
			return
		}

		cardLabel.OwnerUuid = helpers.GetUserUuid(r)
		err = cardLabel.Detach(m.repo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error("Error in detach label")
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info("Detach label success")
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
	params["owner_uuid"] = l.OwnerUuid

	err := repo.ReadManySimple(labels, params)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

func (l *Label) Update(repo Repo) error {
	params := make(map[string]any)
	params["owner_uuid"] = l.OwnerUuid
	params["label_uuid"] = l.LabelUuid

	err := repo.Update(l, params)
	if err != nil {
		return err
	}
	return nil
}

func (l *Label) Delete(repo Repo) error {
	params := make(map[string]any)
	params["owner_uuid"] = l.OwnerUuid
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
	params["owner_uuid"] = c.OwnerUuid
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
	params["owner_uuid"] = c.OwnerUuid

	labels := &[]CardLabel{}
	err := repo.ReadManySimple(labels, params)
	if err != nil {
		return nil, err
	}
	return labels, nil
}
