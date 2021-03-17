package repository

import (
	"fmt"
	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UserPg struct {
	db *sqlx.DB
}

func NewUserPg(db *sqlx.DB) *UserPg {
	return &UserPg{
		db: db,
	}
}

func (r *UserPg) Create(user *domain.User) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var addressId int
	createLocationQuery := fmt.Sprintf(
		`INSERT INTO %s (latitude, longitude)
				VALUES ($1, $2) RETURNING id`, locationsTable)

	locationRow := tx.QueryRow(createLocationQuery, user.Address.Latitude, user.Address.Longitude)
	if err = locationRow.Scan(&addressId); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	createUserQuery := fmt.Sprintf(
		`INSERT INTO %s (name, phone, password_hash, email, address_id)
				VALUES ($1, $2, $3, $4, $5) RETURNING id`, usersTable)

	var userId int
	userRow := tx.QueryRow(createUserQuery, user.Name, user.Phone, user.Password, user.Email, addressId)
	if err = userRow.Scan(&userId); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return userId, tx.Commit()
}

func (r *UserPg) GetByCredentials(phone, password string) (*domain.User, error) {
	user := new(domain.User)
	address := new(domain.Location)

	query := fmt.Sprintf(
		`SELECT u.id, u.name, u.phone, u.password_hash, u.email, l.latitude, l.longitude
				FROM %s AS u JOIN %s AS l ON u.address_id = l.id
				WHERE u.phone = $1 AND u.password_hash = $2`, usersTable, locationsTable)
	row := r.db.QueryRow(query, phone, password)
	err := row.Scan(&user.Id, &user.Name, &user.Phone, &user.Password, &user.Email, &address.Latitude, &address.Longitude)
	user.Address = address

	return user, err
}
