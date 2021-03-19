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
	token, err := s.tokenManager.NewJWT(restaurant.Id, restaurantType, s.accessTokenTTL)

	if err != nil {
		return nil, err
	}

	return &Tokens{AccessToken: token}, nil
}

func (s *RestaurantService) GetAll(clientId int, clientType string) ([]*domain.Restaurant, error) {
	if clientType != userType {
		return nil, errors.New("Forbidden")
	}

	return s.repo.GetAll(clientId)
}

func (s *RestaurantService) GetById(clientId int, clientType string, restaurantId int) (*domain.Restaurant, error) {
	if !(clientType == userType || clientType == restaurantType && restaurantId == clientId) {
		return nil, errors.New("Forbidden")
	}

	return s.repo.GetById(restaurantId)
}

func (s *RestaurantService) GetMenu(clientId int, clientType string, restaurantId int) ([]*domain.MenuItem, error) {
	if !(clientType == userType || clientType == restaurantType && restaurantId == clientId) {
		return nil, errors.New("Forbidden")
	}

	return s.repo.GetMenu(restaurantId)
}

func (s *RestaurantService) SignUp(restaurant *domain.Restaurant, clientType string) (int, error) {
	if clientType != adminType {
		return 0, errors.New("forbidden")
	}
	return s.repo.Create(restaurant)
}
