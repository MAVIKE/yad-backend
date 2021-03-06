package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/MAVIKE/yad-backend/internal/consts"
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

func (s *OrderService) GetAllItems(clientId int, clientType string, orderId int) ([]*domain.OrderItem, error) {
	order, err := s.repo.GetById(orderId)
	if err != nil {
		return nil, err
	}

	if !(clientType == userType && order.UserId == clientId ||
		clientType == restaurantType && order.RestaurantId == clientId ||
		clientType == courierType && order.CourierId == clientId) {
		errMessage := fmt.Sprintf("Forbidden for %s", clientType)
		return nil, errors.New(errMessage)
	}

	items, err := s.repo.GetAllItems(orderId)

	return items, err
}

func (s *OrderService) GetById(clientId int, clientType string, orderId int) (*domain.Order, error) {
	order, err := s.repo.GetById(orderId)

	if err != nil {
		return nil, err
	}

	if !(clientType == userType && order.UserId == clientId ||
		clientType == restaurantType && order.RestaurantId == clientId ||
		clientType == courierType && order.CourierId == clientId) {
		errMessage := fmt.Sprintf("Forbidden for %s", clientType)
		return nil, errors.New(errMessage)
	}

	return order, nil
}

func (s *OrderService) Delete(clientId int, clientType string, orderId int) error {
	order, err := s.repo.GetById(orderId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Order not found")
		}
		return err
	}

	if !(clientType == userType && order.UserId == clientId) {
		return errors.New("Forbidden")
	}

	if order.Status != consts.OrderCreated {
		return errors.New("You can't delete a paid order")
	}

	return s.repo.Delete(orderId)
}

// TODO: обновление общей стоимости заказа лучше сделать в методах добавления, обновления, удаления позиции заказа
func (s *OrderService) Update(clientId int, clientType string, orderId int, input *domain.Order) error {
	order, err := s.repo.GetById(orderId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Order not found")
		}
		return err
	}

	if (input.Status - order.Status) != 1 {
		return errors.New("Invalid new order status")
	}

	switch input.Status {
	case consts.OrderPaid:
		if !(clientType == userType && order.UserId == clientId) {
			errMessage := fmt.Sprintf("Forbidden for %s", clientType)
			return errors.New(errMessage)
		}

		curTime := time.Now()
		input.Paid = &curTime

		courierId, err := s.repo.GetNearestCourierId(order.UserId)
		if err != nil {
			// TODO: свободного курьера может не быть - что делать?
			if err == sql.ErrNoRows {
				return errors.New("Free courier not found")
			}
			return err
		}
		input.CourierId = courierId
	case consts.OrderPreparing, consts.OrderWaitingForCourier:
		if !(clientType == restaurantType && order.RestaurantId == clientId) {
			errMessage := fmt.Sprintf("Forbidden for %s", clientType)
			return errors.New(errMessage)
		}
	case consts.OrderEnRoute, consts.OrderDelivered:
		if !(clientType == courierType && order.CourierId == clientId) {
			errMessage := fmt.Sprintf("Forbidden for %s", clientType)
			return errors.New(errMessage)
		}
	default:
		return errors.New("Order status input error")
	}

	return s.repo.Update(orderId, input)
}

func (s *OrderService) GetActiveRestaurantOrders(clientId int, clientType string, restaurantId int) ([]*domain.Order, error) {
	if !(clientType == userType || clientType == restaurantType && restaurantId == clientId) {
		return nil, errors.New("Forbidden")
	}

	orders, err := s.repo.GetActiveRestaurantOrders(restaurantId)

	return orders, err
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

func (s *OrderService) GetItemById(clientId int, clientType string, orderId, orderItemId int) (*domain.OrderItem, error) {
	order, err := s.repo.GetById(orderId)
	if err != nil {
		return nil, err
	}

	if !(clientType == userType && order.UserId == clientId ||
		clientType == restaurantType && order.RestaurantId == clientId ||
		clientType == courierType && order.CourierId == clientId) {
		errMessage := fmt.Sprintf("Forbidden for %s", clientType)
		return nil, errors.New(errMessage)
	}

	orderItem, err := s.repo.GetItemById(orderItemId)
	if err != nil {
		return nil, err
	}

	if orderItem.OrderId != orderId {
		return nil, errors.New("no such orderItem for this order")
	}

	return orderItem, err
}

func (s *OrderService) UpdateItem(clientId int, clientType string, orderId, orderItemId, menuItemsCount int) error {
	order, err := s.repo.GetById(orderId)
	if err != nil {
		return err
	}

	if !(clientType == userType && order.UserId == clientId) {
		return errors.New("Forbidden")
	}

	orderItem, err := s.repo.GetItemById(orderItemId)
	if err != nil {
		return err
	}

	if orderItem.OrderId != orderId {
		return errors.New("No such orderItem for this order")
	}

	if menuItemsCount < 1 || menuItemsCount > 99 {
		return errors.New("Menu items count must be greater than 0")
	}

	return s.repo.UpdateItem(orderItemId, menuItemsCount)
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

func (s *OrderService) GetActiveCourierOrder(clientId int, clientType string, courierId int) (*domain.Order, error) {
	if !(clientType == courierType && courierId == clientId) {
		errMessage := fmt.Sprintf("Forbidden for %s", clientType)
		return nil, errors.New(errMessage)
	}

	return s.repo.GetActiveCourierOrder(courierId)
}
