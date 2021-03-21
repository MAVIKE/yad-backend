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
	if !(clientType == restaurantType && category.RestaurantId == clientId) {
		return 0, errors.New("Forbidden")
	}

	return s.repo.Create(category)
}

func (s *CategoryService) GetAll(clientId int, clientType string, restaurantId int) ([]*domain.Category, error) {
	if !(clientType == userType || clientType == restaurantType && restaurantId == clientId) {
		return nil, errors.New("forbidden")
	}

	return s.repo.GetAll(restaurantId)
}

func (s *CategoryService) GetById(clientId int, clientType string, restaurantId int, categoryId int) (*domain.Category, error) {
	if !(clientType == userType || clientType == restaurantType && restaurantId == clientId) {
		return nil, errors.New("forbidden")
	}

	category, err := s.repo.GetById(categoryId)
	if err != nil {
		return nil, err
	}

	if category.RestaurantId != restaurantId {
		return nil, errors.New("no such category for this restaurant")
	}
	return category, nil
}
