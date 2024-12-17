package db

func (d *DB) ReadOne(model any, params map[string]any) error {
	result := d.DB.Model(model).Where("user_uuid = ?", params["user_uuid"]).First(model)

	if result.Error != nil {
		return result.Error
	}
	return nil
}
