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
		body, err := helpers.ReadBody(&w, r)
		if err != nil {
			return
		}

		loginUser := user.LoginUser{}
		err = helpers.DeserializeJSON(&w, body, &loginUser)
		if err != nil {
			return
		}

		err = a.validate.Struct(loginUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			for _, err = range err.(validator.ValidationErrors) {
				log.WithField("error", err).Info(types.ValidationError)
			}
			return
		}

		usr := user.User{
			Email: loginUser.Email,
		}

		err = usr.ReadOne(a.repo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.WithField("error", err).Info(types.LoginUserReadFailed)
			return
		}

		if usr.Verified != true {
			w.WriteHeader(http.StatusForbidden)
			log.WithField("user", usr.UserUUID).Warn("Unverified access atempt")
		}

		err = password.ComparePasswords(usr.Password, loginUser.Password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Info(types.PasswordCompareError)
			return
		}

		token, err := jwt.NewJWT(jwt.Secret).Create(usr.UserUUID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.WithField("error", err).Error(types.JwtCreateError)
			return
		}

		w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", token))
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(jwt.TokenResponse{Token: token})
		log.WithField("token", token).Info(types.LoginSuccess)
	}
}

func (a *Auth) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (a *Auth) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := helpers.ReadBody(&w, r)
		if err != nil {
			return
		}
		log.WithFields(log.Fields{"data": string(body)}).Debug(types.ReadBody)

		registerUser := user.RegisterUser{}
		err = helpers.DeserializeJSON(&w, body, &registerUser)
		if err != nil {
			return
		}
		log.WithFields(log.Fields{"data": registerUser}).Debug(types.Deserialized)

		err = a.validate.Struct(registerUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			for _, err = range err.(validator.ValidationErrors) {
				log.WithField("error", err).Error(types.ValidationError)
			}
			return
		}
		log.Debug("Validated")

		registerUser.Password, err = password.HashPassword(registerUser.Password)
		if err != nil {
			http.Error(w, types.PasswordHashError, http.StatusInternalServerError)
			log.WithField("error", err).Error(types.PasswordHashError)
			return
		}
		log.Debug("Password hashed")

		err = registerUser.Create(a.repo)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				w.WriteHeader(http.StatusConflict)
				log.Error(types.UserExists)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			log.WithField("msg", err).Error(types.UserCreateError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		log.Info(types.UserCreateSuccess)
	}
}
