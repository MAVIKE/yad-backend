package service

import (
	"errors"
	"github.com/MAVIKE/yad-backend/internal/domain"
	"time"

	"github.com/MAVIKE/yad-backend/internal/repository"
	"github.com/MAVIKE/yad-backend/pkg/auth"
)

type UserService struct {
	repo           repository.User
	tokenManager   auth.TokenManager
	accessTokenTTL time.Duration
}

func NewUserService(repo repository.User, tokenManager auth.TokenManager, accessTokenTTL time.Duration) *UserService {
	return &UserService{
		repo:           repo,
		tokenManager:   tokenManager,
		accessTokenTTL: accessTokenTTL,
	}
}

func (s *UserService) SignUp(user *domain.User) (int, error) {
	return s.repo.Create(user)
}

func (s *UserService) SignIn(phone, password string) (*Tokens, error) {
	user, err := s.repo.GetByCredentials(phone, password)
	if err != nil {
		return nil, err
	}

	token, err := s.tokenManager.NewJWT(user.Id, userType, s.accessTokenTTL)
	if err != nil {
		return nil, err
	}

	return &Tokens{AccessToken: token}, nil
}

func (s UserService) GetAllOrders(clientId int, clientType string, userId int, activeOrdersFlag bool) ([]*domain.Order, error) {
	if clientType != userType || clientId != userId {
		return nil, errors.New("Forbidden")
	}

	return s.repo.GetAllOrders(clientId, activeOrdersFlag)
}

func (s *UserService) Update(clientId int, clientType string, userId int, input *domain.User) error {

	if !(clientType == userType && userId == clientId) {
		return errors.New("forbidden")
	}

	return s.repo.Update(userId, input)
}

func (s *UserService) GetById(clientId int, clientType string, userId int) (*domain.User, error) {
	if !(clientId == userId || clientType == restaurantType || clientType == courierType) {
		return nil, errors.New("Forbidden")
	}

	// TODO: добавить проверку на то, что у ресторана и курьера есть активный заках с юзером

	return s.repo.GetById(userId)
}
