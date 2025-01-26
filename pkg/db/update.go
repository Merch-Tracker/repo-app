package db

import "fmt"

func (d *DB) Update(model any, params map[string]any) error {
	query := d.DB.Model(model)

	for k, v := range params {
		query = query.Where(fmt.Sprintf("%s = ?", k), v)
	}

	query = query.Updates(model)

	if query.Error != nil {
		return query.Error
	}
	return nil
}

func (d *DB) UpdateNotifications(model any, ids []uint, userUuid string) error {
	const batchSize = 100
	for start := 0; start < len(ids); start += batchSize {
		end := start + batchSize
		if end > len(ids) {
			end = len(ids)
		}

		batch := ids[start:end]

		if err := d.DB.Model(model).
			Where("user_uuid = ?", userUuid).
			Where("id IN ?", batch).
			Update("seen", true).
			Error; err != nil {
			return err
		}
	}
	return nil
}
