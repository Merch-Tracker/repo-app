package images

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"repo-app/pkg/helpers"
	"repo-app/pkg/types"
)

type Repo types.Repo

type ImageHandler struct {
	repo Repo
}

func NewImageHandler(router *http.ServeMux, repo Repo) {
	handler := &ImageHandler{repo: repo}

	err := MigrateImage(repo)
	if err != nil {
		log.WithField("model", "Image").Fatal(types.MerchMigrationError)
	}

	router.HandleFunc("POST /images/{merch_uuid}", handler.Upload())
}

func (i *ImageHandler) Upload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := helpers.ReadBody(&w, r)
		if err != nil {
			return
		}

		merchUuid, err := helpers.GetPathUuid(&w, r, "merch_uuid")

		img := &Image{
			MerchUuid: merchUuid,
		}

		err = helpers.DeserializeJSON(&w, body, img)
		if err != nil {
			return
		}

		err = img.Upload(i.repo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.WithField("error", err).Error("Upload image")
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
