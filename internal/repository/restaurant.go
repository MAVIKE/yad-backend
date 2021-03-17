package repository

import (
	"fmt"
	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type RestaurantPg struct {
	db *sqlx.DB
}

func NewRestaurantPg(db *sqlx.DB) *RestaurantPg {
	return &RestaurantPg{
		db: db,
	}
}

func (r *RestaurantPg) GetByCredentials(phone, password string) (*domain.Restaurant, error) {
	restaurant := new(domain.Restaurant)
	address := new(domain.Location)

	query := fmt.Sprintf(
		`SELECT u.id, u.name, u.phone, u.password_hash, u.email, l.latitude, l.longitude, u.working_status, u.image
 				FROM %s AS u JOIN %s AS l ON u.address_id = l.id
 				WHERE u.phone = $1 AND u.password_hash = $2`, restaurantsTable, locationsTable)
	row := r.db.QueryRow(query, phone, password)
	err := row.Scan(&restaurant.Id, &restaurant.Name, &restaurant.Phone, &restaurant.Password, &restaurant.Email, &address.Latitude, &address.Longitude, &restaurant.WorkingStatus, &restaurant.Image)
	restaurant.Address = address

	return restaurant, err
}
