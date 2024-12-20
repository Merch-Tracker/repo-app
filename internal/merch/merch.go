package merch

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"repo-app/pkg/helpers"
	"repo-app/pkg/types"
)

type Repo types.Repo

type MerchHandler struct {
	repo Repo
}

func NewMerchHandler(router *http.ServeMux, repo Repo) {
	var err error
	handler := &MerchHandler{repo: repo}

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
	router.HandleFunc("PUT /merch", handler.Update())
	router.HandleFunc("DELETE /merch", handler.Delete())
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

		newMerch.OwnerUuid = helpers.GetUUID(r)
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
		readMerch.OwnerUuid = helpers.GetUUID(r)

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
		merch := Merch{OwnerUuid: helpers.GetUUID(r)}

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

		merch := Merch{}
		err = helpers.DeserializeJSON(&w, body, &merch)
		if err != nil {
			return
		}

		merch.OwnerUuid = helpers.GetUUID(r)
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
		merch := Merch{OwnerUuid: helpers.GetUUID(r)}

		err := merch.Delete(m.repo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.WithField("error", err).Error(types.MerchDeleteError)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info(types.MerchDeleteSuccess)
	}
}
