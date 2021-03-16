package service

import (
	"time"

	"github.com/MAVIKE/yad-backend/internal/repository"
	"github.com/MAVIKE/yad-backend/pkg/auth"
)

type UserService struct {
	repo           repository.User
	tokenManager   auth.TokenManager
	accessTokenTTL time.Duration
}

func NewUserService(repo repository.User, tokenManager auth.TokenManager, accessTokenTTL time.Duration) *UserService {
	return &UserService{
		repo:           repo,
		tokenManager:   tokenManager,
		accessTokenTTL: accessTokenTTL,
	}
}

func (s *UserService) SignIn(phone, password string) (*Tokens, error) {
	user, err := s.repo.GetByCredentials(phone, password)
	if err != nil {
		return nil, err
	}

	token, err := s.tokenManager.NewJWT(user.Id, USER_TYPE, s.accessTokenTTL)
	if err != nil {
		return nil, err
	}

	return &Tokens{AccessToken: token}, nil
}

func (s *UserService) ParseToken(accessToken string) (int, string, error) {
	return s.tokenManager.Parse(accessToken)
}
