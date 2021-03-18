package repository

import (
	"fmt"

	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type CategoryPg struct {
	db *sqlx.DB
}

func NewCategoryPg(db *sqlx.DB) *CategoryPg {
	return &CategoryPg{
		db: db,
	}
}

func (r *CategoryPg) Create(category *domain.Category) (int, error) {
	var categoryId int

	query := fmt.Sprintf(
		`INSERT INTO %s (restaurant_id, title)
		VALUES ($1, $2) RETURNING id`, categoriesTable)

	row := r.db.QueryRow(query, category.RestaurantId, category.Title)
	err := row.Scan(&categoryId)

	return categoryId, err
}

func (r *CategoryPg) GetAll(restaurantId int) ([]*domain.Category, error) {
	var categories []*domain.Category

	query := fmt.Sprintf(
		`SELECT c.id, c.restaurant_id, c.title
		FROM %s AS c
		where c.restaurant_id = $1`, categoriesTable)

	err := r.db.Select(&categories, query, restaurantId)

	return categories, err
}
