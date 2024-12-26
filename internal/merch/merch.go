package merch

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"repo-app/pkg/helpers"
	"repo-app/pkg/types"
)

type Repo types.Repo

type MerchHandler struct {
	repo     Repo
	validate *validator.Validate
}

func NewMerchHandler(router *http.ServeMux, repo Repo) {
	var err error
	handler := &MerchHandler{
		repo:     repo,
		validate: validator.New(),
	}

	err = MigrateMerch(repo)
	if err != nil {
		log.WithField("model", "Merch").Fatal(types.MerchMigrationError)
	}

	err = MigrateMerchInfo(repo)
	if err != nil {
		log.WithField("model", "MerchInfo").Fatal(types.MerchMigrationError)
	}

	log.Info(types.MerchMigrationSuccess)

	router.HandleFunc("POST /merch", handler.New())
	router.HandleFunc("GET /merch/", handler.ReadOne())
	router.HandleFunc("GET /merch/all", handler.ReadAll())
	router.HandleFunc("PUT /merch/{merch_uuid}", handler.Update())
	router.HandleFunc("DELETE /merch/{merch_uuid}", handler.Delete())
}

func (m *MerchHandler) New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := helpers.ReadBody(&w, r)
		if err != nil {
			return
		}

		newMerch := Merch{}
		err = helpers.DeserializeJSON(&w, body, &newMerch)
		if err != nil {
			return
		}

		err = m.validate.Struct(newMerch)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			for _, err = range err.(validator.ValidationErrors) {
				log.WithField("error", err).Info(types.ValidationError)
			}
			return
		}

		newMerch.OwnerUuid = helpers.GetUserUuid(r)
		newMerch.MerchUuid = uuid.New()

		err = newMerch.Create(m.repo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.WithField("error", err).Error(types.MerchCreateError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		log.Info(types.MerchCreateSuccess)
	}
}

func (m *MerchHandler) ReadOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		readMerch := Merch{}
		readMerch.OwnerUuid = helpers.GetUserUuid(r)

		err := readMerch.ReadOne(m.repo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.WithField("error", err).Error(types.MerchReadError)
			return
		}

		response, err := helpers.SerializeJSON(&w, readMerch)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.WithField("error", err).Error(types.MerchSerializeError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		log.Info(types.MerchReadSuccess)
	}
}

func (m *MerchHandler) ReadAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		merch := Merch{OwnerUuid: helpers.GetUserUuid(r)}

		allMerch, err := merch.ReadMany(m.repo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.WithField("error", err).Error(types.MerchReadError)
			return
		}

		response, err := helpers.SerializeJSON(&w, allMerch)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.WithField("error", err).Error(types.MerchSerializeError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		log.Info(types.MerchReadSuccess)
	}
}

func (m *MerchHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := helpers.ReadBody(&w, r)
		if err != nil {
			return
		}

		merchUuid, err := helpers.GetPathUuid(&w, r, "merch_uuid")
		if err != nil {
			return
		}

		merch := Merch{}
		err = helpers.DeserializeJSON(&w, body, &merch)
		if err != nil {
			return
		}

		merch.OwnerUuid = helpers.GetUserUuid(r)
		merch.MerchUuid = merchUuid

		err = merch.Update(m.repo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.WithField("error", err).Error(types.MerchUpdateError)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info(types.MerchUpdateSuccess)
	}
}

func (m *MerchHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		merchUuid, err := helpers.GetPathUuid(&w, r, "merch_uuid")
		if err != nil {
			return
		}

		merch := Merch{
			OwnerUuid: helpers.GetUserUuid(r),
			MerchUuid: merchUuid,
		}

		err = merch.Delete(m.repo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.WithField("error", err).Error(types.MerchDeleteError)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info(types.MerchDeleteSuccess)
	}
}
