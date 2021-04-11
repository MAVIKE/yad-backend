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

func (s *RestaurantService) SignUp(restaurant *domain.Restaurant, clientType string) (int, error) {
	if clientType != adminType {
		return 0, errors.New("forbidden")
	}
	return s.repo.Create(restaurant)
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

func (s *RestaurantService) UpdateImage(clientId int, clientType string, restaurantId int, image string) (*domain.Restaurant, error) {
	if clientType != restaurantType || restaurantId != clientId {
		return nil, errors.New("Forbidden")
	}

	if err := s.repo.UpdateImage(restaurantId, image); err != nil {
		return nil, err
	}

	return s.repo.GetById(restaurantId)
}
