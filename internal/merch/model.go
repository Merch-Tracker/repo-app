package merch

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Merch struct {
	gorm.Model `json:"-"`
	MerchUuid  uuid.UUID `gorm:"unique" json:"merch_uuid"`
	UserUuid   uuid.UUID `gorm:"index" json:"-"`
	Name       string    `json:"name" validate:"required,min=1,max=100"`
	Origin     string    `json:"origin"`
}

func (Merch) TableName() string {
	return "merch"
}

type MerchSimple struct {
	CreatedAt time.Time `json:"created_at"`
	MerchUuid string    `json:"merch_uuid"`
	Name      string    `json:"name"`
	Origin    string    `json:"origin"`
}

type Price struct {
	gorm.Model
	MerchUuid uuid.UUID `gorm:"index"`
	Price     uint
}

type NewMerch struct {
	Merch Merch          `json:"merch"`
	Data  map[string]any `json:"data"`
}

type UpdateMerch struct {
	Merch Merch
	Data  any
}

type MerchResponse struct {
	Merch  MerchSimple `json:"merch"`
	Data   any         `json:"data"`
	Labels []uuid.UUID `json:"labels"`
	Prices [2]uint     `json:"prices"`
}

func (n *NewMerch) Create(repo Repo) error {
	n.Data["merch_uuid"] = n.Merch.MerchUuid
	return repo.CreateWithTransaction(&n.Merch, &n.Data, validOrigins[n.Merch.Origin])
}

func (m *Merch) ReadOne(repo Repo) error {
	params := make(map[string]any)
	params["merch_uuid"] = m.MerchUuid
	params["owner_uuid"] = m.UserUuid

	err := repo.ReadOne(m, params)
	if err != nil {
		return err
	}

	fmt.Println(m)

	return nil
}

func (m *MerchResponse) ReadMany(repo Repo, user uuid.UUID) (map[string]MerchResponse, error) {
	params := make(map[string]any, 1)
	params["user_uuid"] = user

	// read user's merch list
	merchList := &[]MerchSimple{}
	err := repo.ReadManySimpleSubmodel(&Merch{}, merchList, params)
	if err != nil {
		return nil, err
	}

	// making response payload
	composed := make(map[string]MerchResponse, len(*merchList))

	// uuids for next db queries
	merchUuids := make(map[string][]string, len(*merchList))

	// insert merch list to response payload
	for _, j := range *merchList {
		if value, exists := merchUuids[j.Origin]; exists {
			merchUuids[j.Origin] = append(value, j.MerchUuid)
		} else {
			merchUuids[j.Origin] = []string{j.MerchUuid}
		}

		composed[j.MerchUuid] = MerchResponse{Merch: j}
	}

	//request additional data
	origin := "surugaya"
	surugaya := &[]Surugaya{}
	err = repo.ReadManyInList(surugaya, merchUuids[origin])
	if err != nil {
		return nil, err
	}

	for _, elem := range *surugaya {
		composed[elem.MerchUuid.String()] = MerchResponse{
			Merch: composed[elem.MerchUuid.String()].Merch,
			Data:  elem}
	}

	origin = "mandarake"
	mandarake := &[]Mandarake{}
	err = repo.ReadManyInList(mandarake, merchUuids[origin])
	if err != nil {
		return nil, err
	}

	for _, elem := range *mandarake {
		composed[elem.MerchUuid.String()] = MerchResponse{
			Merch: composed[elem.MerchUuid.String()].Merch,
			Data:  elem}
	}

	//request labels
	clm := &CardLabel{UserUuid: user}
	cardLabels, err := clm.ReadAll(repo)
	if err != nil {
		return nil, err
	}

	labelList := make(map[string][]uuid.UUID, len(*cardLabels))
	for _, item := range *cardLabels {
		labelList[item.MerchUuid.String()] = append(labelList[item.MerchUuid.String()], item.LabelUuid)
	}

	//request prices
	var allMerch []string
	for _, elem := range *merchList {
		allMerch = append(allMerch, elem.MerchUuid)
	}

	latest := &[]Price{}
	err = repo.ReadPrices(latest, allMerch, 0)
	if err != nil {
		return nil, err
	}

	latestPrices := make(map[string]uint, len(*latest))
	for _, p := range *latest {
		latestPrices[p.MerchUuid.String()] = p.Price
	}

	previous := &[]Price{}
	err = repo.ReadPrices(previous, allMerch, 1)
	if err != nil {
		return nil, err
	}

	previousPrices := make(map[string]uint, len(*previous))
	for _, p := range *previous {
		previousPrices[p.MerchUuid.String()] = p.Price
	}

	for key := range composed {
		composed[key] = MerchResponse{
			Merch:  composed[key].Merch,
			Data:   composed[key].Data,
			Labels: labelList[key],
			Prices: [2]uint{latestPrices[key], previousPrices[key]},
		}
	}

	return composed, nil
}

func (u *UpdateMerch) Update(repo Repo) error {
	params := make(map[string]any, 2)
	params["user_uuid"] = u.Merch.UserUuid
	params["merch_uuid"] = u.Merch.MerchUuid

	err := repo.UpdateWithTransaction(u.Merch, u.Data, params)
	if err != nil {
		return err
	}
	return nil
}

func (u *UpdateMerch) Delete(repo Repo) error {
	params := make(map[string]any, 2)
	params["user_uuid"] = u.Merch.UserUuid
	params["merch_uuid"] = u.Merch.MerchUuid

	err := repo.DeleteWithTransaction(&u.Merch, &u.Data, params)
	if err != nil {
		return err
	}
	return nil
}
