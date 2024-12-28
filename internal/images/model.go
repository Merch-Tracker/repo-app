package images

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	MerchUuid uuid.UUID `json:"MerchUuid"`
	Data      []byte    `json:"Data"`
}

func MigrateImage(repo Repo) error {
	err := repo.Migrate(Image{})
	if err != nil {
		return err
	}
	return nil
}

func (i *Image) Upload(repo Repo) error {
	err := repo.Create(i)
	if err != nil {
		return err
	}
	return nil
}

func (i *Image) Fetch(repo Repo) error {
	params := make(map[string]any)
	params["merch_uuid"] = i.MerchUuid

	err := repo.ReadOne(i, params)
	if err != nil {
		return err
	}
	return nil
}
