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

var validOrigins = map[string]any{
	"surugaya":  &Surugaya{},
	"mandarake": &Mandarake{},
}

type MerchHandler struct {
	repo         Repo
	validate     *validator.Validate
	validOrigins []string
}

func NewMerchHandler(router *http.ServeMux, repo Repo) {
	var err error
	handler := &MerchHandler{
		repo:     repo,
		validate: validator.New(),
	}

	err = migrateMerch(repo)
	if err != nil {
		log.Fatal(merchTableError)
	}

	err = migratePrices(repo)
	if err != nil {
		log.Fatal(pricesTableError)
	}

	err = migrateLabels(repo)
	if err != nil {
		log.Fatal(labelsTableError)
	}

	err = migrateCardLabels(repo)
	if err != nil {
		log.Fatal(cardLabelsTableError)
	}

	err = migrateOriginSurugaya(repo)
	if err != nil {
		log.Fatal(originSurugayaError)
	}

	err = migrateOriginMandarake(repo)
	if err != nil {
		log.Fatal(originMandarakeError)
	}

	log.Debug(migrationsSuccess)

	router.HandleFunc("POST /merch", handler.New())
	router.HandleFunc("GET /merch/all", handler.ReadAll())
	router.HandleFunc("PUT /merch/", handler.Update())
	router.HandleFunc("DELETE /merch/", handler.Delete())

	router.HandleFunc("POST /label", handler.NewLabel())
	router.HandleFunc("GET /label", handler.GetLabels())
	router.HandleFunc("PUT /label/{label_uuid}", handler.UpdateLabel())
	router.HandleFunc("DELETE /label/{label_uuid}", handler.DeleteLabel())

	router.HandleFunc("POST /alabel", handler.AttachLabel())
	router.HandleFunc("POST /dlabel", handler.DetachLabel())

	router.HandleFunc("GET /prices/{merch_uuid}", handler.GetPriceHistory())
	router.HandleFunc("GET /charts", handler.GetAllPrices())
}

func (m *MerchHandler) New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := helpers.ReadBody(w, r)
		if err != nil {
			return
		}

		var newMerch NewMerch
		err = helpers.DeserializeJSON(w, body, &newMerch)
		if err != nil {
			return
		}

		if !validOrigin(newMerch.Merch.Origin) {
			http.Error(w, unknownOriginError, http.StatusBadRequest)
			log.Warn(unknownOriginError)
			return
		}

		err = m.validate.Struct(newMerch.Merch)
		if err != nil {
			http.Error(w, newMerchValidationError, http.StatusBadRequest)
			for _, err = range err.(validator.ValidationErrors) {
				log.WithField(errMsg, err).Info(newMerchValidationError)
			}
			return
		}

		newMerch.Merch.UserUuid = helpers.GetUserUuid(r)
		newMerch.Merch.MerchUuid = uuid.New()

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

func (m *MerchHandler) ReadAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		model := MerchResponse{}
		allMerch, err := model.ReadMany(m.repo, helpers.GetUserUuid(r))
		if err != nil {
			http.Error(w, merchReadError, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(merchReadError)
			return
		}

		response, err := helpers.SerializeJSON(w, allMerch)
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
		body, err := helpers.ReadBody(w, r)
		if err != nil {
			return
		}

		rec := &NewMerch{}
		rec.Merch.UserUuid = helpers.GetUserUuid(r)
		err = helpers.DeserializeJSON(w, body, rec)
		if err != nil {
			return
		}

		upd := &UpdateMerch{
			Merch: rec.Merch,
		}

		switch rec.Merch.Origin {
		case "surugaya":
			upd.Data = &Surugaya{
				MerchUuid:      rec.Merch.MerchUuid,
				Link:           rec.Data["link"].(string),
				ParseTag:       rec.Data["parse_tag"].(string),
				ParseSubstring: rec.Data["parse_substring"].(string),
				CookieValues:   rec.Data["cookie_values"].(string),
				Separator:      rec.Data["separator"].(string),
			}

		case "mandarake":
			upd.Data = &Mandarake{
				MerchUuid: rec.Merch.MerchUuid,
				Link:      rec.Data["link"].(string),
			}

		default:
			http.Error(w, unknownOriginError, http.StatusBadRequest)
			log.Warn(unknownOriginError)
			return
		}

		err = upd.Update(m.repo)
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
		body, err := helpers.ReadBody(w, r)
		if err != nil {
			return
		}

		rec := &NewMerch{}
		rec.Merch.UserUuid = helpers.GetUserUuid(r)
		err = helpers.DeserializeJSON(w, body, rec)
		if err != nil {
			return
		}

		del := &UpdateMerch{
			Merch: rec.Merch,
		}

		switch rec.Merch.Origin {
		case "surugaya":
			del.Data = &Surugaya{
				MerchUuid: rec.Merch.MerchUuid,
			}

		case "mandarake":
			del.Data = &Mandarake{
				MerchUuid: rec.Merch.MerchUuid,
			}

		default:
			http.Error(w, unknownOriginError, http.StatusBadRequest)
			log.Warn(unknownOriginError)
			return
		}

		err = del.Delete(m.repo)
		if err != nil {
			http.Error(w, merchDeleteError, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(merchDeleteError)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info(merchDeleteSuccess)
	}
}

func validOrigin(origin string) bool {
	if _, ok := validOrigins[origin]; ok {
		return true
	}
	return false
}
