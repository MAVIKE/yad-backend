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

func (r *CategoryPg) GetById(categoryId int) (*domain.Category, error) {
	category := new(domain.Category)

	query := fmt.Sprintf(
		`SELECT id, restaurant_id, title
		FROM %s
		WHERE id = $1`,
		categoriesTable)

	row := r.db.QueryRow(query, categoryId)

	err := row.Scan(&category.Id, &category.RestaurantId, &category.Title)

	return category, err
}

func (r *CategoryPg) GetAllItems(categoryId int) ([]*domain.MenuItem, error) {
	var items []*domain.MenuItem

	query := fmt.Sprintf(
		`SELECT mi.id, mi.restaurant_id, mi.title, mi.image, mi.description, mi.price
		FROM
		%s as mi
		join
		%s as ci 
		on mi.id = ci.menu_item_id
		WHERE ci.category_id = $1`, menuItemsTable, categoryItemsTable)
	err := r.db.Select(&items, query, categoryId)

	return items, err
}
