package service

import (
	"github.com/MAVIKE/yad-backend/internal/repository"
)

type Service struct {
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
