package service

import (
	"time"

	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/MAVIKE/yad-backend/internal/repository"
	"github.com/MAVIKE/yad-backend/pkg/auth"
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

type Restaurant interface {
	GetAll(userId int, clientType string) ([]*domain.Restaurant, error)
	GetById(userId int, clientType string, restaurantId int) (*domain.Restaurant, error)
}

type Courier interface {
	SignIn(phone, password string) (*Tokens, error)
}

type Service struct {
	Admin
	User
	Courier
	Restaurant
}

type Deps struct {
	Repos          *repository.Repository
	TokenManager   auth.TokenManager
	AccessTokenTTL time.Duration
}

func NewService(deps Deps) *Service {
	return &Service{
		Admin:      NewAdminService(deps.Repos.Admin, deps.TokenManager, deps.AccessTokenTTL),
		User:       NewUserService(deps.Repos.User, deps.TokenManager, deps.AccessTokenTTL),
		Courier:    NewCourierService(deps.Repos.Courier, deps.TokenManager, deps.AccessTokenTTL),
		Restaurant: NewRestaurantService(deps.Repos.Restaurant),
	}
}
