package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"repo-app/pkg/types"
)

type Repo types.Repo

type User struct {
	gorm.Model
	UserUUID uuid.UUID
	Username string `gorm:"unique"`
	Password string
	Email    string `gorm:"unique"`
}

func Migrate(repo Repo) error {
	err := repo.Migrate(User{})
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Create(repo Repo) error {
	var err error

	u.UserUUID = uuid.New()
	u.Password, err = hashPassword(u.Password)
	if err != nil {
		return err
	}

	err = repo.Create(u)
	if err != nil {
		return err
	}
	return nil
}
