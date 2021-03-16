package repository

import (
	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Admin interface {
	GetByCredentials(name, password string) (*domain.Admin, error)
}

type User interface {
	GetByCredentials(phone, password string) (*domain.User, error)
}

type Courier interface {
	GetByCredentials(phone, password string) (*domain.Courier, error)
}

type Restaurant interface {
	GetAll(userId int) ([]*domain.Restaurant, error)
}

type Repository struct {
	Admin
	User
	Courier
	Restaurant
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Admin:      NewAdminPg(db),
		User:       NewUserPg(db),
		Courier:    NewCourierPg(db),
		Restaurant: NewRestaurantPg(db),
	}
}
