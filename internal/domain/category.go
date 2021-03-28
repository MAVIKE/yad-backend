package domain

type Category struct {
	Id           int    `json:"id" db:"id"`
	RestaurantId int    `json:"restaurant_id" db:"restaurant_id"`
	Title        string `json:"title" db:"title"`
}
