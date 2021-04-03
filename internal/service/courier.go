package service

import (
	"errors"
	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/MAVIKE/yad-backend/internal/repository"
	"github.com/MAVIKE/yad-backend/pkg/auth"
	"time"
)

type CourierService struct {
	repo           repository.Courier
	tokenManager   auth.TokenManager
	accessTokenTTL time.Duration
}

func NewCourierService(repo repository.Courier, tokenManager auth.TokenManager, accessTokenTTL time.Duration) *CourierService {
	return &CourierService{
		repo:           repo,
		tokenManager:   tokenManager,
		accessTokenTTL: accessTokenTTL,
	}
}

func (s *CourierService) SignUp(courier *domain.Courier, clientType string) (int, error) {
	if clientType != adminType {
		return 0, errors.New("forbidden")
	}
	return s.repo.Create(courier)
}

func (s *CourierService) SignIn(phone, password string) (*Tokens, error) {
	courier, err := s.repo.GetByCredentials(phone, password)
	if err != nil {
		return nil, err
	}

	token, err := s.tokenManager.NewJWT(courier.Id, courierType, s.accessTokenTTL)
	if err != nil {
		return nil, err
	}

	return &Tokens{AccessToken: token}, nil
}

func (s *CourierService) GetById(clientId int, clientType string, courierId int) (*domain.Courier, error) {
	if !(clientType == userType || clientType == restaurantType || clientId == courierId) {
		return nil, errors.New("Forbidden")
	}

	return s.repo.GetById(courierId)
}

func (s *CourierService) Update(clientId int, clientType string, courierId int, input *domain.Courier) error {
	if !(clientType == adminType || (clientType == courierType && courierId == clientId)) {
		return errors.New("forbidden")
	}

	if input.WorkingStatus != working && input.WorkingStatus != waiting && input.WorkingStatus != unable {
		return errors.New("working_status input error")
	}

	courier, err := s.repo.GetById(courierId)
	if err != nil {
		return err
	}

	diff := courier.WorkingStatus - input.WorkingStatus
	if diff == 2 || diff == -2 {
		return errors.New("jump over states")
	}

	return s.repo.Update(courierId, input)

}
