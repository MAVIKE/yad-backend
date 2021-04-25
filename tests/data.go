package tests

import "github.com/MAVIKE/yad-backend/internal/domain"

var (
	orders = []domain.Order{
		{
			Id:            1,
			UserId:        1,
			RestaurantId:  1,
			CourierId:     0,
			DeliveryPrice: 100,
			TotalPrice:    900,
			Status:        0,
			Paid:          nil,
		},
		{
			Id:            2,
			UserId:        1,
			RestaurantId:  2,
			CourierId:     1,
			DeliveryPrice: 200,
			TotalPrice:    650,
			Status:        5,
			Paid:          nil,
		},
	}
)
