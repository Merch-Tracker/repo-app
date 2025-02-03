package db

import "fmt"

func (d *DB) ReadOne(model any, params map[string]any) error {
	query := d.DB.Model(model)

	for k, v := range params {
		query = query.Where(fmt.Sprintf("%s = ?", k), v)
	}

	query = query.First(model)
	if query.Error != nil {
		return query.Error
	}

	return nil
}

func (d *DB) ReadOnePayload(model any, payload any, params map[string]any) error {
	query := d.DB.Model(model)

	for k, v := range params {
		query = query.Where(fmt.Sprintf("%s = ?", k), v)
	}

	query = query.First(payload)
	if query.Error != nil {
		return query.Error
	}

	return nil
}

func (d *DB) ReadPrices(model any, list []string, offset int) error {
	return d.DB.Table("prices AS p").
		Select(fmt.Sprintf(`p.id, p.merch_uuid,
		(SELECT price FROM prices AS pi WHERE pi.merch_uuid = p.merch_uuid
		ORDER BY id DESC LIMIT 1 OFFSET %d)`, offset)).
		Where("merch_uuid IN (?)", list).Limit(len(list)).Find(model).Error
}

func (d *DB) Read(model, payload any) error {
	// unused, keep for later
	query := d.DB.Model(model).Find(payload)
	err := query.Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) ReadManySimple(model any, params map[string]any) error {
	query := d.DB.Model(model)

	for k, v := range params {
		query = query.Where(fmt.Sprintf("%s = ?", k), v)
	}

	query = query.Find(model)
	err := query.Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) ReadManySimpleSubmodel(model any, submodel any, params map[string]any) error {
	query := d.DB.Model(model)

	for k, v := range params {
		query = query.Where(fmt.Sprintf("%s = ?", k), v)
	}

	return query.Find(submodel).Error
}

func (d *DB) ReadManyInList(model any, list []string) error {
	return d.DB.Model(model).Where("merch_uuid IN (?)", list).Find(model).Error
}

func (d *DB) ReadManySubmodel(model any, payload any, params map[string]any) error {
	query := d.DB.Model(model)

	for k, v := range params {
		switch k {
		case "days":
			query = query.Where("created_at >= ?", v)

		default:
			query = query.Where(fmt.Sprintf("%s = ?", k), v)
		}
	}

	query.Find(payload)
	err := query.Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) ReadCharts(payload any, params map[string]any) error {
	query := d.DB.Table("prices AS p").
		Select("me.name, p.merch_uuid, json_agg(json_build_object("+
			"'created_at', p.created_at, 'price', p.price) ORDER BY p.created_at) AS prices").
		Joins("JOIN merch AS me ON p.merch_uuid = me.merch_uuid").
		Where("me.user_uuid = ?", params["user_uuid"]).
		Where("p.created_at >= ?", params["days"]).
		Where("p.deleted_at IS NULL").
		Group("p.merch_uuid, me.name").
		Scan(payload)

	err := query.Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) ReadRaw(sql string, payload any) error {
	query := d.DB.Raw(sql).Find(payload)

	err := query.Error
	if err != nil {
		return err
	}
	return nil
}
