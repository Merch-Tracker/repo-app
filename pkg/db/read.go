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

func (d *DB) ReadMany(payload any, params map[string]any) error {
	query := d.DB.Table("merches AS m").
		Select(`m.name, m.link, m.merch_uuid, m.owner_uuid, m.created_at, m.updated_at,
		m.parse_tag, m.parse_substring, m.cookie_values, m.separator,
		(SELECT price FROM merch_infos mi1 WHERE mi1.merch_uuid = m.merch_uuid ORDER BY mi1.id DESC LIMIT 1) AS new_price,
		(SELECT price FROM merch_infos mi2 WHERE mi2.merch_uuid = m.merch_uuid ORDER BY mi2.id DESC OFFSET 1 LIMIT 1) AS old_price`).
		Where("m.deleted_at IS NULL")

	for k, v := range params {
		query = query.Where(fmt.Sprintf("%s = ?", k), v)
	}

	query = query.Find(payload)
	err := query.Error

	if err != nil {
		return err
	}

	return nil
}

func (d *DB) Read(model, payload any) error {
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
	query := d.DB.Table("merch_infos AS mi").
		Select("me.name, me.link, mi.merch_uuid, json_agg(json_build_object("+
			"'created_at', mi.created_at, 'price', mi.price) ORDER BY mi.created_at) AS prices").
		Joins("JOIN merches AS me ON mi.merch_uuid = me.merch_uuid").
		Where("me.owner_uuid = ?", params["owner_uuid"]).
		Where("mi.created_at >= ?", params["days"]).
		Where("mi.deleted_at IS NULL").
		Group("mi.merch_uuid, me.name, me.link").
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
