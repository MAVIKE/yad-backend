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

func (r *CategoryPg) GetAll(restaurantId int) ([]*domain.Category, error) {
	var categories []*domain.Category
	fmt.Println("test")
	query := fmt.Sprintf(
		`SELECT c.id, c.restaurant_id, c.title
		FROM %s AS c
		where c.restaurant_id = $1`, categoriesTable)

	err := r.db.Select(&categories, query, restaurantId)
	fmt.Println("test")
	return categories, err
}
