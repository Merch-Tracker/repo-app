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

func (d *DB) UpdateWithTransaction(model1 any, model2 any, params map[string]any) error {
	tx := d.DB.Begin()
	err := tx.Exec(`SET TRANSACTION ISOLATION LEVEL READ COMMITTED`).Error
	if err != nil {
		return tx.Rollback().Error
	}

	err = tx.
		Where("merch_uuid = ?", params["merch_uuid"]).
		Where("user_uuid = ?", params["user_uuid"]).
		Updates(model1).
		Error
	if err != nil {
		return tx.Rollback().Error
	}

	err = tx.
		Where("merch_uuid = ?", params["merch_uuid"]).
		Updates(model2).
		Error
	if err != nil {
		return tx.Rollback().Error
	}

	return tx.Commit().Error
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
