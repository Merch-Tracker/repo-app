package db

import (
	"fmt"
	"time"
)

func (d *DB) Delete(model any, params map[string]any) error {
	query := d.DB.Model(model)

	for k, v := range params {
		query = query.Where(fmt.Sprintf("%s = ?", k), v)
	}

	query = query.Delete(model)

	if query.Error != nil {
		return query.Error
	}
	return nil
}

func (d *DB) DeleteWithTransaction(model1 any, model2 any, params map[string]any) error {
	tx := d.DB.Begin()
	err := tx.Exec(`SET TRANSACTION ISOLATION LEVEL READ COMMITTED`).Error
	if err != nil {
		return tx.Rollback().Error
	}

	err = tx.
		Where("merch_uuid = ?", params["merch_uuid"]).
		Where("user_uuid = ?", params["user_uuid"]).
		Delete(model1).
		Error
	if err != nil {
		return tx.Rollback().Error
	}

	err = tx.
		Model(model2).
		Where("merch_uuid = ?", params["merch_uuid"]).
		Update("deleted_at", time.Now()).
		Error
	if err != nil {
		return tx.Rollback().Error
	}

	return tx.Commit().Error
}
