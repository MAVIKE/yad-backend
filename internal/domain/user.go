package domain

type User struct {
	Id       int       `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	Phone    string    `json:"phone" db: "phone"`
	Password string    `json:"password" db:"password_hash"`
	Email    string    `json:"email" db:"email""`
	Address  *Location `json:"location" db:"location"`
}
