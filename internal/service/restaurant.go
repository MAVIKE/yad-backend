package service

import (
	"github.com/MAVIKE/yad-backend/internal/repository"
	"github.com/MAVIKE/yad-backend/pkg/auth"
	"time"
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
