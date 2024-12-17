package db

func (d *DB) Delete(model any, params map[string]any) error {
	result := d.DB.Where("user_uuid = ?", params["user_uuid"]).Delete(model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
