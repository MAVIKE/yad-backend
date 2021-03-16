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

func (s *RestaurantService) GetAll(userId int, client_type string) ([]*domain.Restaurant, error) {
	if client_type != USER_TYPE {
		return nil, errors.New("Forbidden")
	}

	restaurants, err := s.repo.GetAll(userId)
	if err != nil {
		return nil, err
	}

	return restaurants, nil
}
