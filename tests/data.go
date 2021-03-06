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
		{
			Id:            4,
			UserId:        3,
			RestaurantId:  1,
			CourierId:     3,
			DeliveryPrice: 100,
			TotalPrice:    700,
			Status:        3,
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
	restaurants = []domain.Restaurant{
		{
			Id:       1,
			Name:     "Restaurant1",
			Phone:    "71234567891",
			Password: "password",
			WorkingStatus: 1,
			Address: &domain.Location{
				Latitude:  52,
				Longitude: 85,
			},
			Image: "img/image1.jpg",
		},
		{
			Id:       2,
			Name:     "Restaurant2",
			Phone:    "71234567892",
			Password: "password",
			WorkingStatus: 1,
			Address: &domain.Location{
				Latitude:  55,
				Longitude: 85,
			},
			Image: "img/image1.jpg",
		},
		{
			Id:       3,
			Name:     "Restaurant2",
			Phone:    "71234567893",
			Password: "password",
			WorkingStatus: 2,
			Address: &domain.Location{
				Latitude:  56,
				Longitude: 87,
			},
			Image: "img/image1.jpg",
		},
  	}
	couriers = []domain.Courier{
		{
			Id:       1,
			Name:     "courier1",
			Phone:    "71234567891",
			Password: "password",
			Email:    "test1@mail.ru",
			Address: &domain.Location{
				Latitude:  52,
				Longitude: 89,
			},
			WorkingStatus: 0,
		},
	}
	menuItems = []domain.MenuItem{
		{
			Id:           1,
			RestaurantId: 1,
			Title:        "Title1",
			Image:        "img/image1.jpg",
			Description:  "description1",
			Price:        100,
		},
		{
			Id:           2,
			RestaurantId: 1,
			Title:        "Title2",
			Image:        "img/image1.jpg",
			Description:  "description2",
			Price:        200,
		},
		{
			Id:           3,
			RestaurantId: 1,
			Title:        "Title3",
			Image:        "img/image1.jpg",
			Description:  "description3",
			Price:        300,
		},
	}
)
