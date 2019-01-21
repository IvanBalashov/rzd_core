package users_gateway

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //?
	"rzd/app/entity"
)

type PgUsers struct {
	DB *sqlx.DB
}

func NewPostgres(db *sqlx.DB) PgUsers {
	return PgUsers{DB: db}
}

func (u *PgUsers) Create(user entity.User) error {
	/*query := `INSERT INTO PgUsers (full_name, nick, train_ids, user_notify) VALUES ($1, $2, $3, $4)`
	_, err := d.DB.Exec(query, user.GetArgs()) // FIXME: check result from query
	if err != nil {
		return err
	}*/
	return nil
}

func (u *PgUsers) Read(offset, limit int) ([]entity.User, error) {
	/*	PgUsers := []entity.User{}
		user := entity.User{}
		query := `SELECT * FROM PgUsers offset $1 limit $2`
		rows, err := d.DB.Queryx(query, offset, limit)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			if err := rows.StructScan(&user); err != nil {
				return nil, err
			}
			PgUsers = append(PgUsers, user)
		}
		return PgUsers, nil*/
	panic("IMPLIMENT ME!")
}

func (u *PgUsers) Update(user entity.User) error {
	panic("IMPLIMENT ME!")
	return nil
}

func (u *PgUsers) Delete(user entity.User) error {
	panic("IMPLIMENT ME!")
	return nil
}
