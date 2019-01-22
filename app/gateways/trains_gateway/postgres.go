package trains_gateway

import (
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //?
	"rzd/app/entity"
)

/*
Ok, i'm write here about correct struct for trains. Coz main think is several pods with `rzd_core`, need development unique
entities for all trains. Need use unique identifier for every `train_entity`.
*/
type PgTrains struct {
	DB *sqlx.DB
}

func NewPostgres(db *sqlx.DB) PgTrains {
	return PgTrains{DB: db}
}

func (t *PgTrains) Create(train entity.Train) error {
	insert := sq.Insert("trains").
		Columns(
			"number",
			"type",
			"brand",
			"route0",
			"route1",
			"trTime0",
			"station",
			"station1",
			"date0",
			"time0",
			"class",
			"seatsCount",
			"price").
		Values(train).
		PlaceholderFormat(sq.Dollar)
	query, args, err := insert.ToSql()
	res, err := t.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	if count, err := res.RowsAffected(); err != nil {
		return err
	} else {
		if count <= 0 {
			return errors.New(fmt.Sprintf("PG:Gateways->Trains_Gateway->Create: Got %s rows affected\n", count))
		}
		return nil
	}
}

func (t *PgTrains) ReadOne() (entity.Train, error) {
	train := entity.Train{}
	query := `SELECT * FROM PgTrains WHERE id = %s`
	rows, err := t.DB.Queryx(fmt.Sprintf(query))
	if err != nil {
		return entity.Train{}, err
	}

	for rows.Next() {
		if err := rows.StructScan(&train); err != nil {
			return entity.Train{}, err
		}
	}
	return train, nil
}

//TODO: think about more effective method for selecting.
func (t *PgTrains) ReadMany(ids []int) ([]entity.Train, error) {
	PgTrains := []entity.Train{}
	train := entity.Train{}
	for _, val := range ids {
		query := `SELECT * FROM PgTrains WHERE id = %s`
		rows, err := t.DB.Queryx(fmt.Sprintf(query, val))
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			if err := rows.StructScan(&train); err != nil {
				return nil, err
			}
			PgTrains = append(PgTrains, train)
		}
	}
	return PgTrains, nil
}

func (t *PgTrains) Update(train entity.Train) error {
	query := `UPDATE PgTrains SET (train_name = %s) where id = %s`
	res, err := t.DB.Exec(fmt.Sprintf(query))
	if err != nil {
		return err
	}
	if count, err := res.RowsAffected(); err != nil {
		return err
	} else {
		if count <= 0 {
			return errors.New(fmt.Sprintf("PG:Gateways->Trains_Gateway->Update: Got %s rows affected\n", count))
		}
		return nil
	}
}

func (t *PgTrains) Delete(train entity.Train) error {
	query := `DELETE FROM trains `
	res, err := t.DB.Exec(query)
	if err != nil {
		return err
	}
	if count, err := res.RowsAffected(); err != nil {
		return err
	} else {
		if count <= 0 {
			return errors.New(fmt.Sprintf("PG:Gateways->Trains_Gateway->Delete: Got %s rows affected\n", count))
		}
		return nil
	}
}
