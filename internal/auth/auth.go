package auth

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"net/http"
	"repo-app/internal/user"
	"repo-app/pkg/helpers"
	"repo-app/pkg/jwt"
	"repo-app/pkg/password"
	"repo-app/pkg/types"
	"strings"
)

type Repo types.Repo

type Auth struct {
	repo     Repo
	validate *validator.Validate
}

func NewAuthHandler(router *http.ServeMux, repo Repo) {
	handler := &Auth{
		repo:     repo,
		validate: validator.New(),
	}

	router.HandleFunc("POST /register", handler.Register())
	router.HandleFunc("POST /login", handler.Login())
	router.HandleFunc("POST /logout", handler.Logout())
}

func (a *Auth) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := helpers.ReadBody(w, r)
		if err != nil {
			return
		}

		loginUser := user.LoginUser{}
		err = helpers.DeserializeJSON(w, body, &loginUser)
		if err != nil {
			return
		}

		err = a.validate.Struct(loginUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			for _, err = range err.(validator.ValidationErrors) {
				log.WithField(errMsg, err).Info(loginValidationError)
			}
			return
		}

		usr := user.User{
			Email: loginUser.Email,
		}

		err = usr.ReadOne(a.repo)
		if err != nil {
			http.Error(w, loginReadUserError, http.StatusBadRequest)
			log.WithField(errMsg, err).Info(loginReadUserError)
			return
		}

		if usr.Verified != true {
			http.Error(w, loginUnverified, http.StatusForbidden)
			log.WithField(respMsg, usr.UserUUID).Warn(loginUnverified)
		}

		err = password.ComparePasswords(usr.Password, loginUser.Password)
		if err != nil {
			http.Error(w, password.PasswordCompareError, http.StatusBadRequest)
			log.Info(password.PasswordCompareError)
			return
		}

		token, err := jwt.NewJWT(jwt.Secret).Create(usr.UserUUID)
		if err != nil {
			http.Error(w, jwt.TokenCreateError, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(jwt.TokenCreateError)
			return
		}

		w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", token))
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(jwt.TokenResponse{Token: token})
		if err != nil {
			http.Error(w, loginResponseEncodeError, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(loginResponseEncodeError)
			return
		}

		log.WithField(respMsg, token).Info(loginSuccess)
	}
}

func (a *Auth) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (a *Auth) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := helpers.ReadBody(w, r)
		if err != nil {
			return
		}

		registerUser := user.RegisterUser{}
		err = helpers.DeserializeJSON(w, body, &registerUser)
		if err != nil {
			return
		}

		err = a.validate.Struct(registerUser)
		if err != nil {
			http.Error(w, registerValidationError, http.StatusBadRequest)
			for _, err = range err.(validator.ValidationErrors) {
				log.WithField(errMsg, err).Error(registerValidationError)
			}
			return
		}

		registerUser.Password, err = password.HashPassword(registerUser.Password)
		if err != nil {
			http.Error(w, password.PasswordHashError, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(password.PasswordHashError)
			return
		}

		err = registerUser.Create(a.repo)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				w.WriteHeader(http.StatusConflict)
				log.Error(registerUserExists)
				return
			}

			http.Error(w, registerUserError, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(registerUserError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		log.Info(registerSuccess)
	}
}
