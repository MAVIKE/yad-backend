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

func (r *OrderPg) GetAllItems(orderId int) ([]*domain.OrderItem, error) {
	var items []*domain.OrderItem

	query := fmt.Sprintf(`SELECT * FROM %s AS oi WHERE oi.order_id = $1`, orderItemsTable)
	err := r.db.Select(&items, query, orderId)

	return items, err
}

func (r *OrderPg) GetById(orderId int) (*domain.Order, error) {
	order := new(domain.Order)

	query := fmt.Sprintf(
		`SELECT o.id, user_id, restaurant_id, COALESCE(courier_id, 0) AS courier_id,
			delivery_price, total_price, status, paid 
		FROM %s AS o WHERE o.id = $1`, ordersTable)
	err := r.db.Get(order, query, orderId)

	return order, err
}

func (r *OrderPg) GetActiveRestaurantOrders(restaurantId int) ([]*domain.Order, error) {
	var orders []*domain.Order

	query := fmt.Sprintf(
		`SELECT * FROM %s 
		WHERE restaurant_id = $1 AND status = %d`, ordersTable, OrderPaid)
	err := r.db.Select(&orders, query, restaurantId)

	return orders, err
}

func (r *OrderPg) CreateItem(orderItem *domain.OrderItem) (int, error) {
	var orderItemId int

	query := fmt.Sprintf(
		`INSERT INTO %s (order_id, menu_item_id, count)
		VALUES ($1, $2, $3) RETURNING id`, orderItemsTable)

	row := r.db.QueryRow(query, orderItem.OrderId, orderItem.MenuItemId, orderItem.Count)
	err := row.Scan(&orderItemId)

	return orderItemId, err
}

func (r *OrderPg) GetItemById(orderItemId int) (*domain.OrderItem, error) {
	item := new(domain.OrderItem)

	query := fmt.Sprintf(`SELECT * FROM %s AS i WHERE i.id = $1`, orderItemsTable)
	err := r.db.Get(item, query, orderItemId)

	return item, err
}

func (r *OrderPg) DeleteItem(orderId int, orderItemId int) error {
	query := fmt.Sprintf(`DELETE FROM %s AS i WHERE i.order_id = $1 AND i.id = $2`, orderItemsTable)
	_, err := r.db.Exec(query, orderId, orderItemId)

	return err
}

func (r *OrderPg) UpdateItem(orderItemId, menuItemsCount int) error {
	query := fmt.Sprintf(`UPDATE %s SET count = $1 WHERE id = $2`, orderItemsTable)
	_, err := r.db.Exec(query, menuItemsCount, orderItemId)

	return err
}
