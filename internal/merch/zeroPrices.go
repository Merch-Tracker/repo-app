package merch

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"repo-app/pkg/helpers"
	"time"
)

type ZeroPrice struct {
	CreatedAt time.Time `json:"created_at"`
	MerchUuid string    `json:"merch_uuid"`
}

//controllers

func (m *MerchHandler) getZeroPrices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		owner := helpers.GetUserUuid(r)

		z := ZeroPrice{}
		zeroPrices, err := z.getZeroPrices(m.repo, owner)
		if err != nil {
			http.Error(w, urlParseError, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(zeroPricesError)
			return
		}

		response, err := helpers.SerializeJSON(w, zeroPrices)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
		log.Info(zeroPricesSuccess)
	}
}

func (m *MerchHandler) deleteZeroPrices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := helpers.ReadBody(w, r)
		if err != nil {
			return
		}

		fmt.Println(string(body))

		if body == nil {
			w.WriteHeader(http.StatusNoContent)
			log.Info("Body is empty")
			return
		}

		toDelete := &[]ZeroPrice{}
		err = helpers.DeserializeJSON(w, body, toDelete)
		if err != nil {
			return
		}

		z := ZeroPrice{}
		err = z.deleteZeroPrices(m.repo, toDelete)
		if err != nil {
			http.Error(w, zeroPricesDeleteError, http.StatusInternalServerError)
			log.WithField(errMsg, err).Error(zeroPricesDeleteError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		log.Info(zeroPricesDeleteSuccess)
	}
}

//services

func (z *ZeroPrice) getZeroPrices(repo Repo, userUuid uuid.UUID) (*[]ZeroPrice, error) {
	prices := &[]ZeroPrice{}

	err := repo.ReadRaw(fmt.Sprintf(`
	WITH ranked_prices AS (
		SELECT 
			id,
			created_at,
			updated_at,
			deleted_at,
			merch_uuid,
			price,
			LAG(price) OVER (PARTITION BY merch_uuid ORDER BY created_at) AS prev_price,
			LEAD(price) OVER (PARTITION BY merch_uuid ORDER BY created_at) AS next_price
		FROM 
			prices
	)
	SELECT 
		r.id,
		r.created_at,
		r.merch_uuid
	FROM 
		ranked_prices AS r
	JOIN merch AS m ON m.merch_uuid = r.merch_uuid
	WHERE
		m.user_uuid = '%s'
		AND price = 0 
		AND prev_price > 0 
		AND next_price > 0
		AND r.deleted_at IS NULL
		;
	`, userUuid), prices)
	if err != nil {
		return nil, err
	}
	return prices, nil
}

func (z *ZeroPrice) deleteZeroPrices(repo Repo, list *[]ZeroPrice) error {
	var err error
	params := make(map[string]any, 2)

	for _, item := range *list {
		params["created_at"] = item.CreatedAt
		params["merch_uuid"] = item.MerchUuid
		params["price"] = 0 //unnecessary condition

		err = repo.Delete(&Price{}, params)
		if err != nil {
			return err
		}
	}

	return nil
}
