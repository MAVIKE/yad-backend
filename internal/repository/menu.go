package repository

import (
	"fmt"
	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type MenuItemPg struct {
	db *sqlx.DB
}

func NewMenuItem(db *sqlx.DB) *MenuItemPg {
	return &MenuItemPg{
		db: db,
	}
}

func (r *MenuItemPg) GetById(menuItemId int) (*domain.MenuItem, error) {
	menuItem := new(domain.MenuItem)

	query := fmt.Sprintf(
		`SELECT m.id, m.restaurant_id, m.title, m.image, m.description
		FROM %s AS m
		WHERE m.id = $1`,
		menuItemsTable)

	row := r.db.QueryRow(query, menuItemId)

	err := row.Scan(&menuItem.Id, &menuItem.RestaurantId, &menuItem.Title, &menuItem.Image,
		&menuItem.Description)

	return menuItem, err
}
