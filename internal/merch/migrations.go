package merch

func migrateMerch(repo Repo) error {
	return repo.Migrate(Merch{})
}

func migratePrices(repo Repo) error {
	return repo.Migrate(Price{})
}

func migrateLabels(repo Repo) error {
	err := repo.Migrate(Label{})
	if err != nil {
		return err
	}
	return nil
}

func migrateCardLabels(repo Repo) error {
	err := repo.Migrate(CardLabel{})
	if err != nil {
		return err
	}
	return nil
}

// Origins

func migrateOriginSurugaya(repo Repo) error {
	return repo.Migrate(Surugaya{})
}

func migrateOriginMandarake(repo Repo) error {
	return repo.Migrate(Mandarake{})
}
