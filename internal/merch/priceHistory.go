package merch

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"repo-app/pkg/helpers"
	"strconv"
	"time"
)

type ChartPoint struct {
	CreatedAt time.Time `json:"date"`
	Price     uint      `json:"price"`
}

func (m *MerchHandler) GetPriceHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		merchUuid, err := helpers.GetPathUuid(&w, r, "merch_uuid")
		if err != nil {
			return
		}

		u, err := url.Parse(r.URL.String())
		if err != nil {
			log.WithField("error", err).Error("Parsing URL query params")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		qParams := u.Query()
		count, err := strconv.Atoi(qParams.Get("count"))
		if err != nil {
			log.WithField("error", err).Error("Parsing query params")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		point := ChartPoint{}
		prices, err := point.GetPriceHistory(m.repo, merchUuid, count)
		if err != nil {
			log.WithField("error", err).Error("Getting price history")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var filteredPrices []ChartPoint

		if prices != nil {
			filteredPrices = append(filteredPrices, (*prices)[0])
			for i := 1; i < len(*prices); i++ {
				if (*prices)[i].Price != (*prices)[i-1].Price {
					filteredPrices = append(filteredPrices, (*prices)[i])
				}
			}
		}

		response, err := helpers.SerializeJSON(&w, filteredPrices)
		if err != nil {
			log.WithField("error", err).Error("Serializing price history")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
		log.Info("Successfully fetched price history")
	}
}

func (c *ChartPoint) GetPriceHistory(repo Repo, merchUuid uuid.UUID, count int) (*[]ChartPoint, error) {
	prices := &[]ChartPoint{}
	params := make(map[string]any)
	params["merch_uuid"] = merchUuid

	if count == 0 || count <= 30 {
		params["days"] = time.Now().Add(-time.Duration(count) * 24 * time.Hour)
	}

	err := repo.ReadManySubmodel(MerchInfo{}, prices, params)
	if err != nil {
		return nil, err
	}

	return prices, nil
}
