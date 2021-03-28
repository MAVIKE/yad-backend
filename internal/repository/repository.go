package repository

import (
	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Admin interface {
	GetByCredentials(name, password string) (*domain.Admin, error)
}

type User interface {
	Create(user *domain.User) (int, error)
	GetByCredentials(phone, password string) (*domain.User, error)
}

type Courier interface {
	Create(courier *domain.Courier) (int, error)
	GetByCredentials(phone, password string) (*domain.Courier, error)
	GetById(courierId int) (*domain.Courier, error)
}

type Restaurant interface {
	GetByCredentials(phone, password string) (*domain.Restaurant, error)
	GetAll(userId int) ([]*domain.Restaurant, error)
	GetById(restaurantId int) (*domain.Restaurant, error)
	GetMenu(restaurantId int) ([]*domain.MenuItem, error)
	Create(restaurant *domain.Restaurant) (int, error)
}

type Category interface {
	GetAll(restaurantId int) ([]*domain.Category, error)
	Create(category *domain.Category) (int, error)
	GetById(categoryId int) (*domain.Category, error)
	GetAllItems(categoryId int) ([]*domain.MenuItem, error)
}

type Order interface {
	Create(order *domain.Order) (int, error)
	GetAllItems(orderId int) ([]*domain.OrderItem, error)
	GetById(orderId int) (*domain.Order, error)
	CreateItem(orderItem *domain.OrderItem) (int, error)
	GetItemById(orderItemId int) (*domain.OrderItem, error)
	UpdateItem(orderItemId, menuItemsCount int) error
	DeleteItem(orderItemId int, orderId int) error
}

type MenuItem interface {
	GetById(menuItemId int) (*domain.MenuItem, error)
}

type Repository struct {
	Admin
	User
	Courier
	Restaurant
	Category
	Order
	MenuItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Admin:      NewAdminPg(db),
		User:       NewUserPg(db),
		Courier:    NewCourierPg(db),
		Restaurant: NewRestaurantPg(db),
		Category:   NewCategoryPg(db),
		Order:      NewOrderPg(db),
		MenuItem:   NewMenuItem(db),
	}
}
