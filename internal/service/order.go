package service

import (
	"errors"
	"fmt"

	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/MAVIKE/yad-backend/internal/repository"
)

type OrderService struct {
	repo repository.Order
}

func NewOrderService(repo repository.Order) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (s *OrderService) Create(clientId int, clientType string, order *domain.Order) (int, error) {
	if clientType != userType {
		return 0, errors.New("Forbidden")
	}

	order.UserId = clientId
	// TODO: установить статус
	// TODO: вычислить и установить стоимость доставки

	return s.repo.Create(order)
}

func (s *OrderService) CreateItem(clientId int, clientType string, orderItem *domain.OrderItem) (int, error) {
	if clientType != userType {
		return 0, errors.New("Forbidden")
	}

	if orderItem.Count < 1 || orderItem.Count > 99 {
		return 0, errors.New("Menu items count must be greater than 0")
	}

	return s.repo.CreateItem(orderItem)
}

func (s *OrderService) GetItemById(clientId int, clientType string, orderId int) (*domain.OrderItem, error) {
	// TODO: проверка прав
	return s.repo.GetItemById(orderId)
}

func (s *OrderService) DeleteItem(clientId int, clientType string, orderId int, orderItemId int) error {
	order, err := s.repo.GetById(orderId)
	if err != nil {
		return err
	}

	if !(clientType == userType && order.UserId == clientId) {
		errMessage := fmt.Sprintf("Forbidden for %s", clientType)
		return errors.New(errMessage)
	}

	return s.repo.DeleteItem(orderId, orderItemId)
}
