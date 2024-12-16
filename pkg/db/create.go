package db

func (d *DB) Migrate(payload any) error {
	err := d.DB.AutoMigrate(payload)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) Create(payload any) error {
	err := d.DB.Create(payload).Error
	if err != nil {
		return err
	}
	return nil
}
