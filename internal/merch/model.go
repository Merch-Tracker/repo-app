package merch

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Merch struct {
	gorm.Model
	MerchUuid uuid.UUID
	OwnerUuid uuid.UUID
	Name      string `validate:"required,min=1,max=100"`
	Link      string
}

type MerchInfo struct {
	gorm.Model
	MerchUuid uuid.UUID
	Price     uint
}

func MigrateMerch(repo Repo) error {
	err := repo.Migrate(Merch{})
	if err != nil {
		return err
	}
	return nil
}

func MigrateMerchInfo(repo Repo) error {
	err := repo.Migrate(MerchInfo{})
	if err != nil {
		return err
	}
	return nil
}

func (m *Merch) Create(repo Repo) error {
	err := repo.Create(m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Merch) ReadOne(repo Repo) error {
	params := make(map[string]any)
	params["merch_uuid"] = m.MerchUuid
	params["owner_uuid"] = m.OwnerUuid

	err := repo.ReadOne(m, params)
	if err != nil {
		return err
	}
	return nil
}

func (m *Merch) ReadMany(repo Repo) (*[]Merch, error) {
	params := make(map[string]any)
	params["owner_uuid"] = m.OwnerUuid

	allMerch := &[]Merch{}

	err := repo.ReadMany(allMerch, params)
	if err != nil {
		return nil, err
	}
	return allMerch, nil
}

func (m *Merch) Update(repo Repo) error {
	params := make(map[string]any)
	params["owner_uuid"] = m.OwnerUuid

	err := repo.Update(m, params)
	if err != nil {
		return err
	}
	return nil
}

func (m *Merch) Delete(repo Repo) error {
	params := make(map[string]any)
	params["owner_uuid"] = m.OwnerUuid

	err := repo.Delete(m, params)
	if err != nil {
		return err
	}
	return nil
}
