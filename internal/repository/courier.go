package repository

import (
	"fmt"
	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type CourierPg struct {
	db *sqlx.DB
}

func NewCourierPg(db *sqlx.DB) *CourierPg {
	return &CourierPg{
		db: db,
	}
}

func (r *CourierPg) GetByCredentials(phone, password string) (*domain.Courier, error) {
	courier := new(domain.Courier)
	address := new(domain.Location)

	query := fmt.Sprintf(
		`SELECT u.id, u.name, u.phone, u.password_hash, u.email, l.latitude, l.longitude, u.working_status
 				FROM %s AS u JOIN %s AS l ON u.address_id = l.id
 				WHERE u.phone = $1 AND u.password_hash = $2`, couriersTable, locationsTable)
	row := r.db.QueryRow(query, phone, password)
	err := row.Scan(&courier.Id, &courier.Name, &courier.Phone, &courier.Password, &courier.Email, &address.Latitude, &address.Longitude, &courier.WorkingStatus)
	courier.Address = address

	return courier, err
}
