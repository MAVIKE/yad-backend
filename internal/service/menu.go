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

func (s *MenuItemService) GetById(clientId int, clientType string, menuItemId int) (*domain.MenuItem, error) {
	if !(clientType == userType || clientType == restaurantType) {
		return nil, errors.New("Forbidden")
	}

	return s.repo.GetById(menuItemId)
}
