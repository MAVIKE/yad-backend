package repository

import (
	"database/sql"
	"errors"
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

func (r *CategoryPg) DeleteCategory(restaurantId int, categoryId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	id := 0
	query := fmt.Sprintf(`SELECT id FROM %s WHERE restaurant_id = $1 AND id = $2`, categoriesTable)
	row := r.db.QueryRow(query, restaurantId, categoryId)
	err = row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			_ = tx.Rollback()
			return errors.New("category does not belong to this restaurant")
		} else {
			_ = tx.Rollback()
			return err
		}
	}

	query = fmt.Sprintf(`DELETE FROM %s WHERE category_id = $1`, categoryItemsTable)
	_, err = r.db.Exec(query, categoryId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	query = fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, categoriesTable)
	_, err = r.db.Exec(query, categoryId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *CategoryPg) UpdateCategory(restaurantId int, categoryId int, input *domain.Category) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	id := 0
	query := fmt.Sprintf(`SELECT id FROM %s WHERE restaurant_id = $1 AND id = $2`, categoriesTable)
	row := r.db.QueryRow(query, restaurantId, categoryId)
	err = row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			_ = tx.Rollback()
			return errors.New("category does not belong to this restaurant")
		}

		_ = tx.Rollback()
		return err
	}

	if input.Title != "" {
		query = fmt.Sprintf(`UPDATE %s SET title = $1 WHERE id = $2`, categoriesTable)
		_, err = r.db.Exec(query, input.Title, categoryId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
