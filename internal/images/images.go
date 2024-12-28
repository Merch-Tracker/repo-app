package images

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"repo-app/pkg/helpers"
	"repo-app/pkg/types"
	"time"
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
	router.HandleFunc("GET /images/{merch_uuid}", handler.SendImage())
}

func (i *ImageHandler) Upload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.WithField("error", err).Error("Error receiving image")
			return
		}

		file, _, err := r.FormFile("Data")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.WithField("error", err).Error("Error receiving image")
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.WithField("error", err).Error("Error receiving image")
			return
		}

		merchUuid, err := helpers.GetPathUuid(&w, r, "merch_uuid")
		if err != nil {
			return
		}

		img := &Image{
			MerchUuid: merchUuid,
			Data:      data,
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

func (i *ImageHandler) SendImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		merchUuid, err := helpers.GetPathUuid(&w, r, "merch_uuid")
		if err != nil {
			return
		}

		img := &Image{MerchUuid: merchUuid}
		err = img.Fetch(i.repo)
		if err.Error() == "record not found" {
			w.WriteHeader(http.StatusNoContent)
			log.Debug("Requested image doesn't exist")
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.WithField("error", err).Error("Error getting image from DB")
			return
		}

		mimeData := http.DetectContentType(img.Data)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", mimeData)
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Header().Set("Expires", time.Now().Add(24*time.Hour).Format(http.TimeFormat))
		w.Write(img.Data)
		log.WithFields(log.Fields{
			"MerchUuid": merchUuid,
			"length":    len(img.Data),
		}).Info("Image fetched")
	}
}
