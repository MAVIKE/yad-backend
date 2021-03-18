package domain

type Category struct {
	Id           int    `json:"id" db:"id"`
	RestaurantId int    `json:"restaurantId" db:"restaurantId"`
	Title        string `json:"title" db:"title"`
}
