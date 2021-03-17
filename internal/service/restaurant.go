package service

import (
	"errors"

	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/MAVIKE/yad-backend/internal/repository"
)

type RestaurantService struct {
	repo repository.Restaurant
}

func NewRestaurantService(repo repository.Restaurant) *RestaurantService {
	return &RestaurantService{repo: repo}
}

func (s *RestaurantService) GetAll(clientId int, clientType string) ([]*domain.Restaurant, error) {
	if clientType != USER_TYPE {
		return nil, errors.New("Forbidden")
	}

	restaurants, err := s.repo.GetAll(clientId)
	if err != nil {
		return nil, err
	}

	return restaurants, nil
}

func (s *RestaurantService) GetById(clientId int, clientType string, restaurantId int) (*domain.Restaurant, error) {
	if clientType != USER_TYPE || clientType == RESTAURANT_TYPE && restaurantId != clientId {
		return nil, errors.New("Forbidden")
	}

	restaurant, err := s.repo.GetById(restaurantId)
	if err != nil {
		return nil, err
	}

	return restaurant, nil
}
