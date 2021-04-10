package repository

import (
	"database/sql"
	"fmt"
	"github.com/MAVIKE/yad-backend/internal/consts"
	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/jmoiron/sqlx"
	"strings"
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

func (r *UserPg) GetAllOrders(userId int, activeOrdersFlag bool) ([]*domain.Order, error) {
	var orders []*domain.Order

	var query string
	var rows *sql.Rows
	var err error

	if activeOrdersFlag {
		query = fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1 and status BETWEEN $2 AND $3`, ordersTable)
		rows, err = r.db.Query(query, userId, consts.OrderPaid, consts.OrderEnRoute)
	} else {
		query = fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1`, ordersTable)
		rows, err = r.db.Query(query, userId)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		order := &domain.Order{}

		err := rows.Scan(&order.Id, &order.UserId, &order.RestaurantId,
			&order.CourierId, &order.DeliveryPrice,
			&order.TotalPrice, &order.Status, &order.Paid)

		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, err
}

func (r *UserPg) Update(userId int, input *domain.User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

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

	if input.Email != "" {
		setValues = append(setValues, fmt.Sprintf("email=$%d", argId))
		args = append(args, input.Email)
		argId++
	}

	if input.Address.Latitude != 0 && input.Address.Longitude != 0 {
		query := fmt.Sprintf(`UPDATE %s as l
								SET latitude = $1, longitude = $2
								FROM %s as c
								WHERE c.id = $3 AND c.address_id = l.id`, locationsTable, usersTable)
		_, err := r.db.Exec(query, input.Address.Latitude, input.Address.Longitude, userId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id=$%d`,
		usersTable, setQuery, argId)
	args = append(args, userId)

	if _, err = tx.Exec(query, args...); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *UserPg) GetById(userId int) (*domain.User, error) {
	user := new(domain.User)
	location := new(domain.Location)

	query := fmt.Sprintf(
		`SELECT u.id, u.name, u.phone, u.email, 
			l.latitude, l.longitude 
		FROM %s AS u
			INNER JOIN %s AS l ON u.address_id = l.id
		WHERE u.id = $1`,
		usersTable, locationsTable)

	row := r.db.QueryRow(query, userId)

	err := row.Scan(&user.Id, &user.Name, &user.Phone, &user.Email,
		&location.Latitude, &location.Longitude)
	user.Address = location

	return user, err
}
