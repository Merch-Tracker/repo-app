package merch

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Merch struct {
	gorm.Model
	MerchUuid      uuid.UUID `gorm:"unique"`
	OwnerUuid      uuid.UUID `gorm:"index"`
	Name           string    `json:"name" validate:"required,min=1,max=100"`
	Link           string    `json:"link"`
	ParseTag       string    `json:"parse_tag"`
	ParseSubstring string    `json:"parse_substring"`
	CookieValues   string    `json:"cookie_values"`
	Separator      string    `json:"separator"`
}

type MerchInfo struct {
	gorm.Model
	MerchUuid uuid.UUID `gorm:"foreignkey:MerchUuid;references:MerchUuid"`
	Price     uint
}

type MerchResponse struct {
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	MerchUuid      uuid.UUID
	OwnerUuid      uuid.UUID
	Name           string `json:"name" validate:"required,min=1,max=100"`
	Link           string `json:"link"`
	OldPrice       uint   `json:"old_price"`
	NewPrice       uint   `json:"new_price"`
	ParseTag       string `json:"parse_tag"`
	ParseSubstring string `json:"parse_substring"`
	CookieValues   string `json:"cookie_values"`
	Separator      string `json:"separator"`
}

type MerchWithLabels struct {
	MerchResponse
	Labels []uuid.UUID `json:"labels"`
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

func (m *Merch) ReadMany(repo Repo) (*[]MerchResponse, error) {
	params := make(map[string]any)
	params["owner_uuid"] = m.OwnerUuid

	allMerch := &[]MerchResponse{}

	err := repo.ReadMany(allMerch, params)
	if err != nil {
		return nil, err
	}
	return allMerch, nil
}

func (m *Merch) Update(repo Repo) error {
	params := make(map[string]any)
	params["owner_uuid"] = m.OwnerUuid
	params["merch_uuid"] = m.MerchUuid

	err := repo.Update(m, params)
	if err != nil {
		return err
	}
	return nil
}

func (m *Merch) Delete(repo Repo) error {
	params := make(map[string]any)
	params["owner_uuid"] = m.OwnerUuid
	params["merch_uuid"] = m.MerchUuid

	err := repo.Delete(m, params)
	if err != nil {
		return err
	}
	return nil
}
