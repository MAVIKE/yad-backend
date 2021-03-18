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
		`SELECT u.id, u.name, u.phone, u.password_hash, l.latitude, l.longitude, u.working_status, u.image
 				FROM %s AS u JOIN %s AS l ON u.address_id = l.id
 				WHERE u.phone = $1 AND u.password_hash = $2`, restaurantsTable, locationsTable)
	row := r.db.QueryRow(query, phone, password)
	err := row.Scan(&restaurant.Id, &restaurant.Name, &restaurant.Phone, &restaurant.Password, &address.Latitude, &address.Longitude, &restaurant.WorkingStatus, &restaurant.Image)
	restaurant.Address = address

	return restaurant, err
}

func (r *RestaurantPg) GetAll(userId int) ([]*domain.Restaurant, error) {
	var restaurants []*domain.Restaurant

	query := fmt.Sprintf(
		`SELECT tmp.id, tmp.name, tmp.phone, tmp.working_status, 
			tmp.latitude, tmp.longitude, tmp.image
		FROM
		(
			SELECT r.id, r.name, r.phone, r.working_status, 
				l.latitude, l.longitude, r.image, 
				get_distance(l.latitude, l.longitude, ua.latitude, ua.longitude) AS distance
			FROM %s AS r
				INNER JOIN %s AS l ON r.address_id = l.id,
				(
					SELECT ul.latitude, ul.longitude
						FROM %s AS u
					INNER JOIN %s AS ul ON u.address_id = ul.id
						WHERE u.id = $1
				) AS ua
			ORDER BY distance
		) AS tmp`,
		restaurantsTable, locationsTable, usersTable, locationsTable)

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		restaurant := &domain.Restaurant{}
		location := &domain.Location{}

		err := rows.Scan(&restaurant.Id, &restaurant.Name, &restaurant.Phone,
			&restaurant.WorkingStatus, &location.Latitude,
			&location.Longitude, &restaurant.Image)

		if err != nil {
			return nil, err
		}

		restaurant.Address = location
		restaurants = append(restaurants, restaurant)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return restaurants, err
}

func (r *RestaurantPg) GetById(restaurantId int) (*domain.Restaurant, error) {
	restaurant := new(domain.Restaurant)
	location := new(domain.Location)

	query := fmt.Sprintf(
		`SELECT r.id, r.name, r.phone, r.working_status, 
			l.latitude, l.longitude, r.image 
		FROM %s AS r
			INNER JOIN %s AS l ON r.address_id = l.id
		WHERE r.id = $1`,
		restaurantsTable, locationsTable)

	row := r.db.QueryRow(query, restaurantId)

	err := row.Scan(&restaurant.Id, &restaurant.Name, &restaurant.Phone, &restaurant.WorkingStatus,
		&location.Latitude, &location.Longitude, &restaurant.Image)
	restaurant.Address = location

	return restaurant, err
}

func (r *RestaurantPg) GetMenu(restarauntId int) ([]*domain.MenuItem, error) {
	var items []*domain.MenuItem

	query := fmt.Sprintf(`SELECT * FROM %s AS m WHERE m.restaurant_id = $1`, menuItemsTable)
	if err := r.db.Select(&items, query, restarauntId); err != nil {
		return nil, err
	}

	return items, nil
}
