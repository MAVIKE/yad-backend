package service

import (
	"github.com/MAVIKE/yad-backend/internal/repository"
	"github.com/MAVIKE/yad-backend/pkg/auth"
	"time"
)

type CourierService struct {
	repo           repository.Courier
	tokenManager   auth.TokenManager
	accessTokenTTL time.Duration
}

func NewCourierService(repo repository.Courier, tokenManager auth.TokenManager, accessTokenTTL time.Duration) *CourierService {
	return &CourierService{
		repo:           repo,
		tokenManager:   tokenManager,
		accessTokenTTL: accessTokenTTL,
	}
}

func (s *CourierService) SignIn(phone, password string) (*Tokens, error) {
	courier, err := s.repo.GetByCredentials(phone, password)
	if err != nil {
		return nil, err
	}

	token, err := s.tokenManager.NewJWT(courier.Id, COURIER_TYPE, s.accessTokenTTL)
	if err != nil {
		return nil, err
	}

	return &Tokens{AccessToken: token}, nil
}