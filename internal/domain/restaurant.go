package domain

type Restaurant struct {
	Id            int       `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Phone         string    `json:"phone" db:"phone"`
	Password      string    `json:"password" db:"password_hash"`
	WorkingStatus int       `json:"working_status" db:"working_status"`
	Address       *Location `json:"location" db:"location"`
	Image         string    `json:"image" db:"image"`
}
