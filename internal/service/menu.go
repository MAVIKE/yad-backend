package service

import (
	"errors"

	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/MAVIKE/yad-backend/internal/repository"
)

type MenuItemService struct {
	repo         repository.MenuItem
	categoryRepo repository.Category
}

func NewMenuItemService(repo repository.MenuItem, categoryRepo repository.Category) *MenuItemService {
	return &MenuItemService{
		repo:         repo,
		categoryRepo: categoryRepo,
	}
}

func (s *RestaurantService) GetMenu(clientId int, clientType string, restaurantId int) ([]*domain.MenuItem, error) {
	if !(clientType == userType || clientType == restaurantType && restaurantId == clientId) {
		return nil, errors.New("Forbidden")
	}

	return s.repo.GetMenu(restaurantId)
}

func (s *MenuItemService) GetById(clientId int, clientType string, menuItemId int, restaurantId int) (*domain.MenuItem, error) {
	if !(clientType == userType || clientType == restaurantType && restaurantId == clientId) {
		return nil, errors.New("Forbidden")
	}

	menuItem, err := s.repo.GetById(menuItemId)

	if err != nil {
		return nil, err
	}

	if menuItem.RestaurantId != restaurantId {
		return nil, errors.New("No such menu item for this restaurant")
	}

	return menuItem, nil
}

func (s *MenuItemService) UpdateMenuItem(clientId int, clientType string, restaurantId int, menuItemId int, categoryId int, input *domain.MenuItem) error {
	if !(clientType == restaurantType && restaurantId == clientId) {
		return errors.New("forbidden")
	}

	return s.repo.UpdateMenuItem(restaurantId, menuItemId, categoryId, input)
}

func (s *MenuItemService) Create(clientId int, clientType string, menuItem *domain.MenuItem, categoryId int) (int, error) {
	if !(clientType == restaurantType && menuItem.RestaurantId == clientId) {
		return 0, errors.New("forbidden")
	}

	category, err := s.categoryRepo.GetById(categoryId)
	if clientId != category.RestaurantId || err != nil {
		return 0, errors.New("forbidden")
	}

	return s.repo.Create(menuItem, categoryId)
}
