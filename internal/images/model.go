package images

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	MerchUuid uuid.UUID
	Data      []byte
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
