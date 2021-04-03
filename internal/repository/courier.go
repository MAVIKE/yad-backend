package repository

import (
	"errors"
	"fmt"
	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/jmoiron/sqlx"
	"strings"
)

type CourierPg struct {
	db *sqlx.DB
}

func NewCourierPg(db *sqlx.DB) *CourierPg {
	return &CourierPg{
		db: db,
	}
}

func (r *CourierPg) Create(courier *domain.Courier) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var addressId int
	createLocationQuery := fmt.Sprintf(
		`INSERT INTO %s (latitude, longitude)
 				VALUES ($1, $2) RETURNING id`, locationsTable)

	locationRow := tx.QueryRow(createLocationQuery, courier.Address.Latitude, courier.Address.Longitude)
	if err = locationRow.Scan(&addressId); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	createCourierQuery := fmt.Sprintf(
		`INSERT INTO %s (name, phone, password_hash, email, address_id, working_status)
 				VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, couriersTable)

	var courierId int
	userRow := tx.QueryRow(createCourierQuery, courier.Name, courier.Phone, courier.Password, courier.Email, addressId, courier.WorkingStatus)
	if err = userRow.Scan(&courierId); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return courierId, tx.Commit()
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

func (r *CourierPg) GetById(courierId int) (*domain.Courier, error) {
	courier := new(domain.Courier)
	location := new(domain.Location)

	query := fmt.Sprintf(
		`SELECT c.id, c.name, c.phone, c.email, c.working_status, 
			l.latitude, l.longitude 
		FROM %s AS c
			INNER JOIN %s AS l ON c.address_id = l.id
		WHERE c.id = $1`,
		couriersTable, locationsTable)

	row := r.db.QueryRow(query, courierId)

	err := row.Scan(&courier.Id, &courier.Name, &courier.Phone, &courier.Email, &courier.WorkingStatus,
		&location.Latitude, &location.Longitude)
	courier.Address = location

	return courier, err
}

func (r *CourierPg) Update(courierId int, input *domain.Courier) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

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

	if input.Phone != "" {
		setValues = append(setValues, fmt.Sprintf("phone=$%d", argId))
		args = append(args, input.Phone)
		argId++
	}

	if input.Email != "" {
		setValues = append(setValues, fmt.Sprintf("email=$%d", argId))
		args = append(args, input.Email)
		argId++
	}

	// TODO : добавить поле со статусом и локацией

	if len(setValues) == 0 {
		return nil
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id=$%d`,
		couriersTable, setQuery, argId)
	args = append(args, courierId)

	_, err := r.db.Exec(query, args...)
	print(err)

	return errors.New(query)
}
