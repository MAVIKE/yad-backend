package service

import (
	"errors"

	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/MAVIKE/yad-backend/internal/repository"
)

type MenuItemService struct {
	repo repository.MenuItem
}

func NewMenuItemService(repo repository.MenuItem) *MenuItemService {
	return &MenuItemService{
		repo: repo,
	}
}

func (s *MenuItemService) GetById(clientId int, clientType string, menuItemId int, restaurantId int) (*domain.MenuItem, error) {
	if !(clientType == userType || clientType == restaurantType && restaurantId == clientId) {
		return nil, errors.New("Forbidden")
	}

	menuItem, err := s.repo.GetById(menuItemId)

	if err != nil {
		return  nil, err
	}

	if menuItem.RestaurantId != restaurantId {
		return nil, errors.New("No such menu item for this restaurant")
	}

	return menuItem, nil
}
