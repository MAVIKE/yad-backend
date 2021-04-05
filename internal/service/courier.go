package service

import (
	"errors"
	"fmt"
	"github.com/MAVIKE/yad-backend/internal/consts"
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
	switch input.WorkingStatus {
	case consts.CourierUnable, consts.CourierWaiting, consts.CourierWorking:
		break
	default:
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

	// TODO: проверить, что у курьера нет активного заказа

	if clientType == courierType && courierId == clientId {
		input.Email = ""
		input.Name = ""
		input.Phone = ""
		input.Password = ""
	} else if clientType != adminType {
		return errors.New("forbidden")
	}

	return s.repo.Update(courierId, input)
}

func (s *CourierService) GetActiveOrder(clientId int, clientType string, courierId int) (*domain.Order, error) {
	if clientType != courierType || (clientType == courierType && courierId != clientId) {
		errMessage := fmt.Sprintf("Forbidden for %s", clientType)
		return nil, errors.New(errMessage)
	}

	return s.repo.GetActiveOrder(courierId)
}
