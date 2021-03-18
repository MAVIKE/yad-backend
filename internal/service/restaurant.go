package service

import (
	"errors"
	"time"

	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/MAVIKE/yad-backend/internal/repository"
	"github.com/MAVIKE/yad-backend/pkg/auth"
)

type RestaurantService struct {
	repo           repository.Restaurant
	tokenManager   auth.TokenManager
	accessTokenTTL time.Duration
}

func NewRestaurantService(repo repository.Restaurant, tokenManager auth.TokenManager, accessTokenTTL time.Duration) *RestaurantService {
	return &RestaurantService{
		repo:           repo,
		tokenManager:   tokenManager,
		accessTokenTTL: accessTokenTTL,
	}
}

func (s *RestaurantService) SignIn(phone, password string) (*Tokens, error) {
	restaurant, err := s.repo.GetByCredentials(phone, password)
	if err != nil {
		return nil, err
	}
	token, err := s.tokenManager.NewJWT(restaurant.Id, RESTAURANT_TYPE, s.accessTokenTTL)

	if err != nil {
		return nil, err
	}

	return &Tokens{AccessToken: token}, nil
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
	if !(clientType == USER_TYPE || clientType == RESTAURANT_TYPE && restaurantId == clientId) {
		return nil, errors.New("Forbidden")
	}

	restaurant, err := s.repo.GetById(restaurantId)
	if err != nil {
		return nil, err
	}

	return restaurant, nil
}
