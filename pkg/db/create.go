package db

func (d *DB) Migrate(model any) error {
	return d.DB.AutoMigrate(model)
}

func (d *DB) Create(payload any) error {
	return d.DB.Create(payload).Error
}

func (d *DB) CreateOrRewrite(payload any) error {
	return d.DB.Save(payload).Error
}

func (d *DB) CreateWithTransaction(payload1, payload2 any, origin any) error {
	tx := d.DB.Begin()
	err := tx.Exec(`SET TRANSACTION ISOLATION LEVEL READ COMMITTED`).Error
	if err != nil {
		return tx.Rollback().Error
	}

	err = tx.Create(payload1).Error
	if err != nil {
		return tx.Rollback().Error
	}

	err = tx.Model(origin).Create(payload2).Error
	if err != nil {
		return tx.Rollback().Error
	}

	return tx.Commit().Error
}
