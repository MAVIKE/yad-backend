package service

import (
	"time"

	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/MAVIKE/yad-backend/internal/repository"
	"github.com/MAVIKE/yad-backend/pkg/auth"
)

const (
	adminType      = "admin"
	userType       = "user"
	courierType    = "courier"
	restaurantType = "restaurant"
)

type Tokens struct {
	AccessToken string `json:"token"`
}

type Admin interface {
	SignIn(name, password string) (*Tokens, error)
}

type User interface {
	SignIn(phone, password string) (*Tokens, error)
	SignUp(user *domain.User) (int, error)
}

type Restaurant interface {
	GetAll(clientId int, clientType string) ([]*domain.Restaurant, error)
	GetById(clientId int, clientType string, restaurantId int) (*domain.Restaurant, error)
	SignIn(phone, password string) (*Tokens, error)
	GetMenu(clientId int, clientType string, restaurantId int) ([]*domain.MenuItem, error)
	SignUp(restaurant *domain.Restaurant, clientType string) (int, error)
}

type Courier interface {
	SignIn(phone, password string) (*Tokens, error)
	SignUp(courier *domain.Courier, clientType string) (int, error)
	GetById(clientId int, clientType string, courierId int) (*domain.Courier, error)
}

type Category interface {
	GetAll(clientId int, clientType string, restaurantId int) ([]*domain.Category, error)
	Create(clientId int, clientType string, category *domain.Category) (int, error)
	GetById(clientId int, clientType string, restaurantId int, categoryId int) (*domain.Category, error)
}

type Order interface {
	Create(clientId int, clientType string, order *domain.Order) (int, error)
	GetAllItems(clientId int, clientType string) ([]*domain.OrderItem, error)
	GetById(clientId int, clientType string, orderId int) (*domain.Order, error)
	CreateItem(clientId int, clientType string, orderItem *domain.OrderItem) (int, error)
	GetItemById(clientId int, clientType string, orderId, orderItemId int) (*domain.OrderItem, error)
}

type MenuItem interface {
	GetById(clientId int, clientType string, menuItemId int, restaurantId int) (*domain.MenuItem, error)
}

type Service struct {
	Admin
	User
	Courier
	Restaurant
	Category
	Order
	MenuItem
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
		Restaurant: NewRestaurantService(deps.Repos.Restaurant, deps.TokenManager, deps.AccessTokenTTL),
		Category:   NewCategoryService(deps.Repos.Category),
		Order:      NewOrderService(deps.Repos.Order),
		MenuItem:   NewMenuItemService(deps.Repos.MenuItem),
	}
}
