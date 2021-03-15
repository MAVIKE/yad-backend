package repository

import (
	"fmt"
	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type AdminPg struct {
	db *sqlx.DB
}

func NewAdminPg(db *sqlx.DB) *AdminPg {
	return &AdminPg{
		db: db,
	}
}

func (r *AdminPg) GetByCredentials(name, password string) (*domain.Admin, error) {
	admin := new(domain.Admin)

	query := fmt.Sprintf(`SELECT * FROM %s AS a WHERE a.name = $1 AND a.password_hash = $2`, adminsTable)
	if err := r.db.Get(admin, query, name, password); err != nil {
		return nil, err
	}

	return admin, nil
}
