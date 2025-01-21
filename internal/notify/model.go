package notify

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type NotifyMessage struct {
	gorm.Model
	UserUuid  string `json:"-"`
	MerchUuid string `json:"merch_uuid"`
	Price     uint   `json:"price"`
	Viewed    bool   `json:"seen"`
}

type NotifyMessageResponse struct {
	UserUuid  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	MerchUuid string    `json:"merch_uuid"`
	Price     uint      `json:"price"`
	Viewed    bool      `json:"seen"`
}

type Notifier struct {
	gorm.Model
	UserUuid uuid.UUID `json:"user_uuid"`
	Target   string    `json:"target"`
	Origin   string    `json:"origin"`
}

type UsersList struct {
	UserUuid string `json:"user_uuid"`
	Target   string `json:"target"`
	Origin   string `json:"origin"`
}

type PricesList struct {
	UserUuid  string `json:"user_uuid"`
	MerchUuid string `json:"merch_uuid"`
	Name      string `json:"name"`
	Price     uint   `json:"price"`
}

type NotifierRecord struct {
	Target string `json:"target"`
	Origin string `json:"origin"`
}

type PriceRecords struct {
	MerchUuid string `json:"merch_uuid"`
	Price     uint   `json:"price"`
}

type Response struct {
	Notifiers []NotifierRecord `json:"notifiers"`
	Prices    []PriceRecords   `json:"prices"`
}

func MigrateNotifiers(repo Repo) error {
	err := repo.Migrate(Notifier{})
	if err != nil {
		return err
	}
	return nil
}

func MigrateNotifyMessages(repo Repo) error {
	err := repo.Migrate(NotifyMessage{})
	if err != nil {
		return err
	}
	return nil
}

func (n *UsersList) GetList(repo Repo) (*[]UsersList, error) {
	sql := `
		SELECT u.user_uuid, n.target, n.origin
		FROM users AS u
		JOIN notifiers as n on u.user_uuid = n.user_uuid
		WHERE u.deleted_at IS NULL
		`

	payload := &[]UsersList{}

	err := repo.ReadRaw(sql, payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (p *PricesList) GetList(repo Repo, userList []string) (*[]PricesList, error) {
	sql := fmt.Sprintf(`
	WITH RankedPrices AS ( SELECT mi.merch_uuid, mi.price,
		ROW_NUMBER() OVER (PARTITION BY mi.merch_uuid ORDER BY mi.created_at DESC) AS rn
	FROM merch_infos mi )
	SELECT u.user_uuid, m.merch_uuid, m.name, rp.price
	FROM users u
	JOIN merches m ON u.user_uuid = m.owner_uuid
	JOIN RankedPrices rp ON m.merch_uuid = rp.merch_uuid
	WHERE rp.rn <= 2 AND u.user_uuid IN (%s)
	ORDER BY u.user_uuid, m.merch_uuid, rp.rn`, `'`+strings.Join(userList, ",")+`'`)

	payload := &[]PricesList{}

	err := repo.ReadRaw(sql, payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (n *NotifyMessageResponse) ReadAll(repo Repo) (*[]NotifyMessageResponse, error) {
	sql := fmt.Sprintf(`
		SELECT m.name, nm.created_at, nm.merch_uuid, nm.price
		FROM notify_messages AS nm
		JOIN merches AS m ON m.merch_uuid = nm.merch_uuid
		WHERE nm.user_uuid = '%s'
		AND nm.deleted_at IS NULL
		ORDER BY nm.created_at DESC;
	`, n.UserUuid)

	messages := &[]NotifyMessageResponse{}

	err := repo.ReadRaw(sql, messages)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (n *NotifyMessage) MarkAsRead(repo Repo, list *[]NotifyMessage) error {
	err := repo.Save(list)
	if err != nil {
		return err
	}
	return nil
}
