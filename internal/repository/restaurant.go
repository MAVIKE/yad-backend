package repository

import (
	"fmt"
	"strings"

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

func (r *RestaurantPg) Create(restaurant *domain.Restaurant) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var addressId int
	createLocationQuery := fmt.Sprintf(
		`INSERT INTO %s (latitude, longitude)
 				VALUES ($1, $2) RETURNING id`, locationsTable)

	locationRow := tx.QueryRow(createLocationQuery, restaurant.Address.Latitude, restaurant.Address.Longitude)
	if err = locationRow.Scan(&addressId); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	createRestaurantQuery := fmt.Sprintf(
		`INSERT INTO %s (name, phone, password_hash, address_id, working_status, image)
 				VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, restaurantsTable)

	var restaurantId int
	restaurantRow := tx.QueryRow(createRestaurantQuery, restaurant.Name, restaurant.Phone, restaurant.Password, addressId, restaurant.WorkingStatus, restaurant.Image)
	if err = restaurantRow.Scan(&restaurantId); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return restaurantId, tx.Commit()
}

func (r *RestaurantPg) UpdateImage(restaurantId int, image string) error {
	query := fmt.Sprintf(`UPDATE %s AS r SET image = $1 WHERE r.id = $2`, restaurantsTable)
	_, err := r.db.Exec(query, image, restaurantId)
	return err
}

func (r *RestaurantPg) Update(restaurantId int, input *domain.Restaurant) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.WorkingStatus != 0 {
		setValues = append(setValues, fmt.Sprintf("working_status=$%d", argId))
		args = append(args, input.WorkingStatus)
		argId++
	}

	if input.Name != "" {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, input.Name)
		argId++
	}

	if input.Password != "" {
		setValues = append(setValues, fmt.Sprintf("password_hash=$%d", argId))
		args = append(args, input.Password)
		argId++
	}

	if input.Address.Latitude != 0 && input.Address.Longitude != 0 {
		query := fmt.Sprintf(`UPDATE %s as l
								SET latitude = $1, longitude = $2
								FROM %s as c
								WHERE c.id = $3 AND c.address_id = l.id`, locationsTable, restaurantsTable)
		_, err := r.db.Exec(query, input.Address.Latitude, input.Address.Longitude, restaurantId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id=$%d`,
		restaurantsTable, setQuery, argId)
	args = append(args, restaurantId)

	if _, err = tx.Exec(query, args...); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
