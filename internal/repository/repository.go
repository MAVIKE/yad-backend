package repository

import (
	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Admin interface {
	GetByCredentials(name, password string) (*domain.Admin, error)
}

type User interface {
	Create(user *domain.User) (int, error)
	GetByCredentials(phone, password string) (*domain.User, error)
}

type Courier interface {
	Create(courier *domain.Courier) (int, error)
	GetByCredentials(phone, password string) (*domain.Courier, error)
}

type Restaurant interface {
	GetByCredentials(phone, password string) (*domain.Restaurant, error)
	GetAll(userId int) ([]*domain.Restaurant, error)
	GetById(restarauntId int) (*domain.Restaurant, error)
	GetMenu(restarauntId int) ([]*domain.MenuItem, error)
}

type Category interface {
	GetAll(restaurantId int) ([]*domain.Category, error)
	Create(category *domain.Category) (int, error)
}

type Repository struct {
	Admin
	User
	Courier
	Restaurant
	Category
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Admin:      NewAdminPg(db),
		User:       NewUserPg(db),
		Courier:    NewCourierPg(db),
		Restaurant: NewRestaurantPg(db),
		Category:   NewCategoryPg(db),
	}
}
