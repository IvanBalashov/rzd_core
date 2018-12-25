package trains_gateway

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //?
	"rzd/app/entity"
	"strings"
)

type PgTrains struct {
	DB *sqlx.DB
}

func NewPostgres(db *sqlx.DB) PgTrains {
	return PgTrains{DB: db}
}

func (d *PgTrains) Create(train entity.Train) error {
	query := `INSERT INTO Trains 
    (number, type, brand, route0, route1, trTime0, station, station1, date0, time0, class, seatsCount, price) 
    VALUES ('%s')`
	_, err := d.DB.Exec(fmt.Sprintf(query, strings.Join(train.GetArgs(), "', '")))
	if err != nil {
		return err
	}
	return nil
}

func (d *PgTrains) ReadOne(id int) (entity.Train, error) {
	train := entity.Train{}
	query := `SELECT * FROM PgTrains WHERE id = %s`
	rows, err := d.DB.Queryx(fmt.Sprintf(query, id))
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
func (d *PgTrains) ReadMany(ids []int) ([]entity.Train, error) {
	PgTrains := []entity.Train{}
	train := entity.Train{}
	for _, val := range ids {
		query := `SELECT * FROM PgTrains WHERE id = %s`
		rows, err := d.DB.Queryx(fmt.Sprintf(query, val))
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

func (d *PgTrains) Update(train entity.Train) error {
	query := `UPDATE PgTrains SET (train_name = %s) where id = %s`
	_, err := d.DB.Exec(query, train.GetArgs()) // FIXME: check result from query
	if err != nil {
		return err
	}
	return nil
}

func (d *PgTrains) Delete(train entity.Train) error {
	query := `DELETE FROM trains `
	_, err := d.DB.Exec(query, train.GetArgs()) // FIXME: check result from query
	if err != nil {
		return err
	}
	return nil
}
