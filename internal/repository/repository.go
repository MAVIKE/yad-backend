package repository

import (
	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Admin interface {
	GetByCredentials(name, password string) (*domain.Admin, error)
}

type Repository struct {
	Admin
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Admin: NewAdminPg(db),
	}
}
