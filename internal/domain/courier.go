package domain

type Courier struct {
	Id            int    `json:"id" db:"id"`
	Name          string `json:"name" db:"name"`
	Phone         string `json:"phone" db:"phone"`
	Password      string `json:"password" db:"password_hash"`
	Email         string `json:"email" db:"email"`
	Address       int    `json:"address" db:"address_id"`
	WorkingStatus int    `json:"working_status" db:"working_status"`
}
