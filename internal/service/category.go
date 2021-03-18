package service

import (
	"errors"

	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/MAVIKE/yad-backend/internal/repository"
)

type CategoryService struct {
	repo repository.Category
}

func NewCategoryService(repo repository.Category) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (s *CategoryService) Create(clientId int, clientType string, category *domain.Category) (int, error) {
	if !(clientType == RESTAURANT_TYPE && category.RestaurantId == clientId) {
		return 0, errors.New("Forbidden")
	}

	return s.repo.Create(category)
}

func (s *CategoryService) GetAll(clientId int, clientType string, restaurantId int) ([]*domain.Category, error) {
	if !(clientType == USER_TYPE || clientType == RESTAURANT_TYPE && restaurantId == clientId) {
		return nil, errors.New("Forbidden")
	}

	return s.repo.GetAll(restaurantId)
}
