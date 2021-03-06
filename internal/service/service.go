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
	GetAllOrders(clientId int, clientType string, userId int, activeOrdersFlag bool) ([]*domain.Order, error)
	Update(clientId int, clientType string, userId int, input *domain.User) error
	GetById(clientId int, clientType string, userId int) (*domain.User, error)
}

type Restaurant interface {
	GetAll(clientId int, clientType string) ([]*domain.Restaurant, error)
	GetById(clientId int, clientType string, restaurantId int) (*domain.Restaurant, error)
	SignIn(phone, password string) (*Tokens, error)
	GetMenu(clientId int, clientType string, restaurantId int) ([]*domain.MenuItem, error)
	SignUp(restaurant *domain.Restaurant, clientType string) (int, error)
	UpdateImage(clientId int, clientType string, restaurantId int, image string) (*domain.Restaurant, error)
	Update(clientId int, clientType string, restaurantId int, input *domain.Restaurant) error
}

type Courier interface {
	SignIn(phone, password string) (*Tokens, error)
	SignUp(courier *domain.Courier, clientType string) (int, error)
	GetById(clientId int, clientType string, courierId int) (*domain.Courier, error)
	Update(clientId int, clientType string, courierId int, input *domain.Courier) error
}

type Category interface {
	GetAll(clientId int, clientType string, restaurantId int) ([]*domain.Category, error)
	Create(clientId int, clientType string, category *domain.Category) (int, error)
	GetById(clientId int, clientType string, restaurantId int, categoryId int) (*domain.Category, error)
	GetAllItems(clientId int, clientType string, restaurantId int, categoryId int) ([]*domain.MenuItem, error)
	DeleteCategory(clientId int, clientType string, restaurantId int, categoryId int) error
	UpdateCategory(clientId int, clientType string, restaurantId int, categoryId int, input *domain.Category) error
}

type Order interface {
	Create(clientId int, clientType string, order *domain.Order) (int, error)
	GetById(clientId int, clientType string, orderId int) (*domain.Order, error)
	Delete(clientId int, clientType string, orderId int) error
	Update(clientId int, clientType string, orderId int, status *domain.Order) error
	GetActiveRestaurantOrders(clientId int, clientType string, restaurantId int) ([]*domain.Order, error)
	CreateItem(clientId int, clientType string, orderItem *domain.OrderItem) (int, error)
	GetAllItems(clientId int, clientType string, orderId int) ([]*domain.OrderItem, error)
	GetItemById(clientId int, clientType string, orderId, orderItemId int) (*domain.OrderItem, error)
	UpdateItem(clientId int, clientType string, orderId, orderItemId, menuItemsCount int) error
	DeleteItem(clientId int, clientType string, orderId int, orderItemId int) error
	GetActiveCourierOrder(clientId int, clientType string, courierId int) (*domain.Order, error)
}

type MenuItem interface {
	GetById(clientId int, clientType string, menuItemId int, restaurantId int) (*domain.MenuItem, error)
	UpdateMenuItem(clientId int, clientType string, restaurantId int, menuItemId int, categoryId int, input *domain.MenuItem) error
	Create(clientId int, clientType string, menuItem *domain.MenuItem, categoryId int) (int, error)
	UpdateImage(clientId int, clientType string, restaurantId int, menuItemId int, image string) (*domain.MenuItem, error)
	Delete(clientId int, clientType string, restaurantId int, menuItemId int) error
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
		Courier:    NewCourierService(deps.Repos.Courier, deps.Repos.Order, deps.TokenManager, deps.AccessTokenTTL),
		Restaurant: NewRestaurantService(deps.Repos.Restaurant, deps.TokenManager, deps.AccessTokenTTL),
		Category:   NewCategoryService(deps.Repos.Category),
		Order:      NewOrderService(deps.Repos.Order),
		MenuItem:   NewMenuItemService(deps.Repos.MenuItem, deps.Repos.Category),
	}
}
