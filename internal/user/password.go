package user

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"repo-app/pkg/types"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Error(types.PasswordHashErr)
		return "", err
	}
	return string(bytes), nil
}

func comparePasswords(hashedPassword string, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
