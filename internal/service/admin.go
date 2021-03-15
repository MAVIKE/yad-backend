package service

import (
	"github.com/MAVIKE/yad-backend/internal/repository"
	"github.com/MAVIKE/yad-backend/pkg/auth"
	"time"
)

type AdminService struct {
	repo           repository.Admin
	tokenManager   auth.TokenManager
	accessTokenTTL time.Duration
}

func NewAdminService(repo repository.Admin, tokenManager auth.TokenManager, accessTokenTTL time.Duration) *AdminService {
	return &AdminService{
		repo:           repo,
		tokenManager:   tokenManager,
		accessTokenTTL: accessTokenTTL,
	}
}

func (s *AdminService) SignIn(name, password string) (*Tokens, error) {
	admin, err := s.repo.GetByCredentials(name, password)
	if err != nil {
		return nil, err
	}

	token, err := s.tokenManager.NewJWT(admin.Id, ADMIN_TYPE, s.accessTokenTTL)
	if err != nil {
		return nil, err
	}

	return &Tokens{AccessToken: token}, nil
}
