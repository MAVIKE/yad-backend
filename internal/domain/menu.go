package domain

type MenuItem struct {
	Id           int    `json:"id" db:"id"`
	RestaurantId int    `json:"restaurant_id" db:"restaurant_id"`
	Title        string `json:"title" db:"title"`
	Image        string `json:"image" db:"image"`
	Description  string `json:"description" db:"description"`
	Price        int    `json:"price" db:"price"`
}
