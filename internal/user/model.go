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
	Username string `gorm:"unique" validate:"required"`
	Password string `validate:"required"`
	Email    string `gorm:"unique" validate:"required"`
}

func Migrate(repo Repo) error {
	err := repo.Migrate(User{})
	if err != nil {
		return err
	}
	return nil
}

// Create Creates new user record in repository using RegisterUser payload.
func (r *RegisterUser) Create(repo Repo) error {
	usr := User{
		UserUUID: uuid.New(),
		Username: r.Username,
		Password: r.Password,
		Email:    r.Email,
	}

	err := repo.Create(&usr)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) ReadOne(repo Repo) error {
	params := make(map[string]any)

	switch {
	case u.Email != "":
		params["email"] = u.Email
	case u.UserUUID != uuid.Nil:
		params["user_uuid"] = u.UserUUID

	}

	err := repo.ReadOne(u, params)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) ReadOnePayload(repo Repo, payload any) error {
	params := map[string]any{"user_uuid": u.UserUUID}
	err := repo.ReadOnePayload(u, payload, params)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Update(repo Repo) error {
	params := make(map[string]any)
	params["user_uuid"] = u.UserUUID

	err := repo.Update(u, params)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Delete(repo Repo) error {
	params := make(map[string]any)
	params["user_uuid"] = u.UserUUID

	err := repo.Delete(u, params)
	if err != nil {
		return err
	}
	return nil
}
