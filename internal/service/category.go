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

func (s *CategoryService) GetAll(clientId int, clientType string, restaurantId int) ([]*domain.Category, error) {
	if clientType != USER_TYPE || clientType == RESTAURANT_TYPE && restaurantId != clientId {
		return nil, errors.New("Forbidden")
	}

	categories, err := s.repo.GetAll(restaurantId)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
