package service

import (
	"backend/domain/models"
	"backend/repository"
)

type UserService interface {
	CreateUser(email, name string) (*models.User, error)
	GetUser(id uint) (*models.User, error)
	GetUsers() ([]models.User, error)
	UpdateUserStatus(id uint, isActive bool) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(email, name string) (*models.User, error) {
	user := &models.User{
		Email: email,
		Name:  name,
	}
	err := s.repo.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUser(id uint) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) GetUsers() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *userService) UpdateUserStatus(id uint, isActive bool) error {
	return s.repo.UpdateStatus(id, isActive)
}
