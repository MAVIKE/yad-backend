package domain

import "time"

type Order struct {
	Id            int        `json:"id" db:"id"`
	UserId        int        `json:"user_id" db:"user_id"`
	RestaurantId  int        `json:"restaurant_id" db:"restaurant_id"`
	CourierId     int        `json:"courier_id" db:"courier_id"`
	DeliveryPrice int        `json:"delivery_price" db:"delivery_price"`
	TotalPrice    int        `json:"total_price" db:"total_price"`
	Status        int        `json:"status" db:"status"`
	Paid          *time.Time `json:"paid" db:"paid"`
}
