package merch

import (
	"github.com/google/uuid"
	"time"
)

type Surugaya struct {
	Id             uint      `gorm:"primary_key" json:"-"`
	DeletedAt      time.Time `gorm:"index" json:"-"`
	MerchUuid      uuid.UUID `gorm:"index" json:"-"`
	Link           string    `json:"link"`
	ParseTag       string    `json:"parse_tag"`
	ParseSubstring string    `json:"parse_substring"`
	CookieValues   string    `json:"cookie_values"`
	Separator      string    `json:"separator"`
}

func (Surugaya) TableName() string {
	return "origin_surugaya"
}

type Mandarake struct {
	Id        uint      `gorm:"primary_key" json:"-"`
	DeletedAt time.Time `gorm:"index" json:"-"`
	MerchUuid uuid.UUID `gorm:"index" json:"-"`
	Link      string    `json:"link"`
}

func (Mandarake) TableName() string {
	return "origin_mandarake"
}
