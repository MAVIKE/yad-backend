package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

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

func (r *RestaurantPg) GetMenu(restarauntId int) ([]*domain.MenuItem, error) {
	var items []*domain.MenuItem

	query := fmt.Sprintf(`SELECT * FROM %s AS m WHERE m.restaurant_id = $1`, menuItemsTable)
	if err := r.db.Select(&items, query, restarauntId); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *MenuItemPg) GetById(menuItemId int) (*domain.MenuItem, error) {
	menuItem := new(domain.MenuItem)

	query := fmt.Sprintf(
		`SELECT m.id, m.restaurant_id, m.title, m.image, m.description, m.price
		FROM %s AS m
		WHERE m.id = $1`,
		menuItemsTable)

	row := r.db.QueryRow(query, menuItemId)

	err := row.Scan(&menuItem.Id, &menuItem.RestaurantId, &menuItem.Title, &menuItem.Image,
		&menuItem.Description, &menuItem.Price)

	return menuItem, err
}

func (r *MenuItemPg) UpdateMenuItem(restaurantId int, menuItemId int, categoryId int, input *domain.MenuItem) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	id := 0
	query := fmt.Sprintf(`SELECT id FROM %s WHERE restaurant_id = $1 AND id = $2`, menuItemsTable)
	row := r.db.QueryRow(query, restaurantId, menuItemId)
	err = row.Scan(&id)
	if err == sql.ErrNoRows {
		_ = tx.Rollback()
		return errors.New("menu item does not belong to this restaurant")
	}

	if categoryId != 0 {
		id := 0
		query := fmt.Sprintf(`SELECT id FROM %s WHERE restaurant_id = $1 AND id = $2`, categoriesTable)
		row := r.db.QueryRow(query, restaurantId, categoryId)
		err := row.Scan(&id)
		if err == sql.ErrNoRows {
			_ = tx.Rollback()
			return errors.New("category does not belong to this restaurant")
		}

		query = fmt.Sprintf(`DELETE FROM %s WHERE menu_item_id = $1`, categoryItemsTable)
		_, err = r.db.Exec(query, menuItemId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		query = fmt.Sprintf(`INSERT INTO %s (category_id, menu_item_id) VALUES($1, $2)`, categoryItemsTable)
		_, err = r.db.Exec(query, categoryId, menuItemId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if input.Title != "" {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, input.Title)
		argId++
	}

	if input.Image != "" {
		setValues = append(setValues, fmt.Sprintf("image=$%d", argId))
		args = append(args, input.Image)
		argId++
	}

	if input.Description != "" {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, input.Description)
		argId++
	}

	if input.Price != 0 {
		setValues = append(setValues, fmt.Sprintf("price=$%d", argId))
		args = append(args, input.Price)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query = fmt.Sprintf(`UPDATE %s SET %s WHERE id=$%d`,
		menuItemsTable, setQuery, argId)
	args = append(args, menuItemId)

	if _, err = tx.Exec(query, args...); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *MenuItemPg) Create(menuItem *domain.MenuItem, categoryId int) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var menuItemId int

	query := fmt.Sprintf(
		`INSERT INTO %s (restaurant_id, title, image, description, price)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`, menuItemsTable)
	row := r.db.QueryRow(query, menuItem.RestaurantId, menuItem.Title, menuItem.Image, menuItem.Description,
		menuItem.Price)
	err = row.Scan(&menuItemId)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	var categoryItemsId int
	query = fmt.Sprintf(
		`INSERT INTO %s (category_id, menu_item_id)
		VALUES ($1, $2) RETURNING id`, categoryItemsTable)
	row = r.db.QueryRow(query, categoryId, menuItemId)
	err = row.Scan(&categoryItemsId)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return menuItemId, tx.Commit()
}

func (r *MenuItemPg) UpdateImage(menuItemId int, image string) error {
	query := fmt.Sprintf(`UPDATE %s AS r SET image = $1 WHERE r.id = $2`, menuItemsTable)
	_, err := r.db.Exec(query, image, menuItemId)
	return err
}

func (r *MenuItemPg) DeleteItem(menuItemId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`DELETE FROM %s WHERE menu_item_id = $1`, categoryItemsTable)
	_, err = r.db.Exec(query, menuItemId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	query = fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, menuItemsTable)
	_, err = r.db.Exec(query, menuItemId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
