package repository

import (
	"fmt"

	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type OrderPg struct {
	db *sqlx.DB
}

func NewOrderPg(db *sqlx.DB) *OrderPg {
	return &OrderPg{
		db: db,
	}
}

func (r *OrderPg) Create(order *domain.Order) (int, error) {
	var orderId int

	query := fmt.Sprintf(
		`INSERT INTO %s (user_id, restaurant_id, delivery_price, total_price, status)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`, ordersTable)

	row := r.db.QueryRow(query, order.UserId, order.RestaurantId, order.DeliveryPrice,
		order.TotalPrice, order.Status)
	err := row.Scan(&orderId)

	return orderId, err
}