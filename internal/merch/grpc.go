package merch

import (
	"github.com/google/uuid"
)

type GMerch struct {
	MerchUuid      uuid.UUID
	Link           string
	ParseTag       string
	ParseSubstring string
	CookieValues   string
	Separator      string
}

func (m *GMerch) ReadAll(repo Repo) ([]GMerch, error) {
	var allMerch []GMerch
	sql := `
		SELECT m.merch_uuid, os.link, os.parse_tag, os.parse_substring, os.cookie_values, os.separator
		FROM merch as m 
		JOIN origin_surugaya as os ON m.merch_uuid = os.merch_uuid;   
	`

	err := repo.ReadRaw(sql, &allMerch)
	if err != nil {
		return nil, err
	}
	return allMerch, nil
}
