package db

func (d *DB) Update(model any, params map[string]any) error {
	result := d.DB.Where("user_uuid = ?", params["user_uuid"]).Updates(model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
