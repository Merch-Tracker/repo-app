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
		log.Fatal(merchTableError)
	}

	err = MigrateMerchInfo(repo)
	if err != nil {
		log.Fatal(pricesTableError)
	}

	err = MigrateLabel(repo)
	if err != nil {
		log.Fatal(labelsTableError)
	}

	err = MigrateCardLabel(repo)
	if err != nil {
		log.Fatal(cardLabelsTableError)
	}

	log.Info(migrationsSuccess)

	router.HandleFunc("POST /merch", handler.New())
	router.HandleFunc("GET /merch/", handler.ReadOne())
	router.HandleFunc("GET /merch/all", handler.ReadAll())
	router.HandleFunc("PUT /merch/{merch_uuid}", handler.Update())
	router.HandleFunc("DELETE /merch/{merch_uuid}", handler.Delete())

	router.HandleFunc("POST /label", handler.NewLabel())
	router.HandleFunc("GET /label", handler.GetLabels())
	router.HandleFunc("PUT /label/{label_uuid}", handler.UpdateLabel())
	router.HandleFunc("DELETE /label/{label_uuid}", handler.DeleteLabel())

	router.HandleFunc("POST /alabel", handler.AttachLabel())
	router.HandleFunc("POST /dlabel", handler.DetachLabel())

	router.HandleFunc("GET /prices/{merch_uuid}", handler.GetPriceHistory())
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
			http.Error(w, newMerchValidationError, http.StatusBadRequest)
			for _, err = range err.(validator.ValidationErrors) {
				log.WithField(errMsg, err).Info(newMerchValidationError)
			}
			return
		}

		newMerch.OwnerUuid = helpers.GetUserUuid(r)
		newMerch.MerchUuid = uuid.New()

		err = newMerch.Create(m.repo)
		if err != nil {
			http.Error(w, merchCreateError, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(merchCreateError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		log.Info(merchCreateSuccess)
	}
}

func (m *MerchHandler) ReadOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		readMerch := Merch{}
		readMerch.OwnerUuid = helpers.GetUserUuid(r)

		err := readMerch.ReadOne(m.repo)
		if err != nil {
			http.Error(w, merchReadError, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(merchReadError)
			return
		}

		response, err := helpers.SerializeJSON(&w, readMerch)
		if err != nil {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		log.Info(merchReadSuccess)
	}
}

func (m *MerchHandler) ReadAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		owner := helpers.GetUserUuid(r)
		merch := Merch{OwnerUuid: owner}

		// select all merch
		allMerch, err := merch.ReadMany(m.repo)
		if err != nil {
			http.Error(w, merchReadError, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(merchReadError)
			return
		}

		//select labels for merch
		label := CardLabel{OwnerUuid: owner}
		cardLabels, err := label.ReadAll(m.repo)
		if err != nil {
			http.Error(w, labelsGetAllError, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(labelsGetAllError)
			return
		}

		labelList := make(map[uuid.UUID][]uuid.UUID, len(*cardLabels))
		for _, item := range *cardLabels {
			labelList[item.MerchUuid] = append(labelList[item.MerchUuid], item.LabelUuid)
		}

		//composing
		var composedResponse []MerchWithLabels

		for _, item := range *allMerch {
			composedResponse = append(composedResponse, MerchWithLabels{item, labelList[item.MerchUuid]})
		}

		// composed response
		response, err := helpers.SerializeJSON(&w, composedResponse)
		if err != nil {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		log.WithField(bytesMsg, len(response)).Info(merchReadSuccess)
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
			http.Error(w, merchUpdateError, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(merchUpdateError)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info(merchUpdateSuccess)
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
			http.Error(w, merchDeleteError, http.StatusInternalServerError)
			log.WithField("error", err).Error(merchDeleteError)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info(merchDeleteSuccess)
	}
}
