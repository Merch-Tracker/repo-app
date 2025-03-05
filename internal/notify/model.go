package notify

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type NotifyMessage struct {
	Id        uint   `json:"id" gorm:"primary_key"`
	UserUuid  string `json:"-"`
	MerchUuid string `json:"merch_uuid"`
	PriceId   uint   `json:"price_id" gorm:"unique"`
	Seen      bool   `json:"seen"`
}

type NotifyMessageResponse struct {
	UserUuid  string    `json:"-"`
	NoteId    uint      `json:"note_id" gorm:"column:id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	MerchUuid string    `json:"merch_uuid"`
	Price     uint      `json:"price"`
	Seen      bool      `json:"seen"`
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
	PriceId   uint   `json:"price_id"`
}

type NotifierRecord struct {
	Target string `json:"target"` //link or id
	Origin string `json:"origin"` //platform name
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
		LEFT JOIN notifiers as n on u.user_uuid = n.user_uuid
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
	quotedUserList := make([]string, len(userList))
	for i, userUUID := range userList {
		quotedUserList[i] = "'" + userUUID + "'"
	}

	sql := fmt.Sprintf(`
	WITH RankedPrices AS ( SELECT p.merch_uuid, p.price, p.id AS price_id,
		ROW_NUMBER() OVER (PARTITION BY p.merch_uuid ORDER BY p.created_at DESC) AS rn
	FROM prices AS p )
	SELECT u.user_uuid, m.merch_uuid, m.name, rp.price, price_id
	FROM users u
	JOIN merch AS m ON u.user_uuid = m.user_uuid
	JOIN RankedPrices rp ON m.merch_uuid = rp.merch_uuid
	WHERE rp.rn <= 2 AND u.user_uuid IN (%s)
	AND m.deleted_at IS NULL
	ORDER BY u.user_uuid, m.merch_uuid, rp.rn;`, strings.Join(quotedUserList, ","))

	payload := &[]PricesList{}

	err := repo.ReadRaw(sql, payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (n *NotifyMessageResponse) ReadAll(repo Repo) (*[]NotifyMessageResponse, error) {
	sql := fmt.Sprintf(`
		SELECT nm.id, m.name, mi.created_at, nm.merch_uuid, mi.price, nm.seen
		FROM notify_messages AS nm
		JOIN merch AS m ON m.merch_uuid = nm.merch_uuid
		JOIN prices AS mi ON nm.price_id = mi.id
		WHERE nm.user_uuid = '%s'
		ORDER BY mi.created_at DESC;
	`, n.UserUuid)

	messages := &[]NotifyMessageResponse{}

	err := repo.ReadRaw(sql, messages)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (n *NotifyMessage) MarkAsRead(repo Repo, list []uint, userUuid string) error {
	return repo.UpdateNotifications(n, list, userUuid)
}
