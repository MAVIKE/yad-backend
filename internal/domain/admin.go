package domain

type Admin struct {
	Id       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Password string `json:"password" db:"password_hash"`
}
