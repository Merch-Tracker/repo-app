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
			http.Error(w, urlParseError, http.StatusBadRequest)
			log.WithField(errMsg, err).Error(urlParseError)
			return
		}

		qParams := u.Query()
		count, err := strconv.Atoi(qParams.Get("count"))
		if err != nil {
			http.Error(w, queryParseError, http.StatusBadRequest)
			log.WithField(errMsg, err).Error(queryParseError)
			return
		}

		point := ChartPoint{}
		prices, err := point.GetPriceHistory(m.repo, merchUuid, count)
		if err != nil {
			http.Error(w, getPriceHistoryError, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(getPriceHistoryError)
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
			http.Error(w, serPriceHistory, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(serPriceHistory)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
		log.Info(priceHistoryFetchSuccess)
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
