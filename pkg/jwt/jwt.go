package jwt

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const JWTSecret = "b7dfb67ac0109337fa2f2612c213be7bfc63252e148bf63637bccf6ee4cddc1d"

type JWT struct {
	SecretKey string
}

type Data struct {
	UserID uuid.UUID
}

type TokenResponse struct {
	Token string `json:"token"`
}

func NewJWT(secretKey string) *JWT {
	return &JWT{
		SecretKey: secretKey,
	}
}

func (j *JWT) Create(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
	})

	signedToken, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *JWT) Parse(token string) (*Data, error) {
	parse, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	strId := parse.Claims.(jwt.MapClaims)["userID"].(string)
	id, err := uuid.Parse(strId)
	if err != nil {
		return nil, err
	}

	return &Data{UserID: id}, nil
}
