package service

import (
	"github.com/MAVIKE/yad-backend/internal/repository"
	"github.com/MAVIKE/yad-backend/pkg/auth"
	"time"
)

type Tokens struct {
	AccessToken string `json:"token"`
}

type Admin interface {
	SignIn(name, password string) (*Tokens, error)
}

type User interface {
	SignIn(phone, password string) (*Tokens, error)
}

type Service struct {
	Admin
	User
}

type Deps struct {
	Repos          *repository.Repository
	TokenManager   auth.TokenManager
	AccessTokenTTL time.Duration
}

func NewService(deps Deps) *Service {
	return &Service{
		Admin: NewAdminService(deps.Repos.Admin, deps.TokenManager, deps.AccessTokenTTL),
		User:  NewUserService(deps.Repos.User, deps.TokenManager, deps.AccessTokenTTL),
	}
}
