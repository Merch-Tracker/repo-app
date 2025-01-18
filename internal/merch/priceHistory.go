package merch

import (
	"encoding/json"
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

type ChartsData struct {
	Name      string          `json:"name"`
	Link      string          `json:"link"`
	MerchUuid uuid.UUID       `json:"MerchUuid"`
	Prices    json.RawMessage `json:"prices"`
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
			count = 7 //make it init in handler options
		}

		point := ChartPoint{}
		prices, err := point.GetPriceHistory(m.repo, merchUuid, count)
		if err != nil {
			http.Error(w, getPriceHistoryError, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(getPriceHistoryError)
			return
		}

		if len(*prices) == 0 {
			w.WriteHeader(http.StatusNoContent)
			log.WithField(respMsg, noContentMsg).Info(priceHistoryFetchSuccess)
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
	params["days"] = countDays(count)

	err := repo.ReadManySubmodel(MerchInfo{}, prices, params)
	if err != nil {
		return nil, err
	}

	return prices, nil
}

func (m *MerchHandler) GetAllPrices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		owner := helpers.GetUserUuid(r)

		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, urlParseError, http.StatusBadRequest)
			log.WithField(errMsg, err).Error(urlParseError)
			return
		}

		qParams := u.Query()
		count, err := strconv.Atoi(qParams.Get("count"))
		if err != nil {
			count = 7 //make it init in handler options
		}

		point := ChartsData{}
		chartData, err := point.GetAllPrices(m.repo, owner, count)
		if err != nil {
			http.Error(w, chartsReadDataError, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(chartsReadDataError)
			return
		}

		if chartData == nil {
			w.WriteHeader(http.StatusNoContent)
			log.WithField(respMsg, "No content").Info(chartsReadDataSuccess)
		}

		response, err := helpers.SerializeJSON(&w, chartData)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
		log.WithField(bytesMsg, len(response)).Info(chartsReadDataSuccess)
	}
}

func (c *ChartsData) GetAllPrices(repo Repo, ownerUuid uuid.UUID, count int) (*[]ChartsData, error) {
	prices := &[]ChartsData{}
	params := make(map[string]any)
	params["owner_uuid"] = ownerUuid
	params["days"] = countDays(count)

	err := repo.ReadCharts(prices, params)
	if err != nil {
		return nil, err
	}
	return prices, nil
}

func countDays(count int) time.Time {
	if count > 0 && count <= 30 {
		return time.Now().Add(-time.Duration(count) * 24 * time.Hour)
	} else {
		return time.Now().Add(-time.Duration(7) * 24 * time.Hour)
	}
}
