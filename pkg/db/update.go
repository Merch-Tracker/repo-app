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
