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
	err := repo.Read(Merch{}, &allMerch)
	if err != nil {
		return nil, err
	}
	return allMerch, nil
}
