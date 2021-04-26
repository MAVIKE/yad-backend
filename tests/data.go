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
		{
			Id:            3,
			UserId:        2,
			RestaurantId:  2,
			CourierId:     2,
			DeliveryPrice: 100,
			TotalPrice:    800,
			Status:        1,
			Paid:          nil,
		},
	}

	users = []domain.User{
		{
			Id:       1,
			Name:     "user1",
			Phone:    "71234567890",
			Password: "password",
			Email:    "test1@mail.ru",
			Address: &domain.Location{
				Latitude:  50,
				Longitude: 87,
			},
		},
	}
)
