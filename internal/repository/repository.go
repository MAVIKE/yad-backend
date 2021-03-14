package repository

import (
	_ "github.com/jmoiron/sqlx"
)

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}
