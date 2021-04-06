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

func (r *UserPg) GetAllOrders(userId int) ([]*domain.Order, error) {
	var orders []*domain.Order

	query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1`, ordersTable)

	rows, err := r.db.Query(query, userId)
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
