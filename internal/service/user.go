package service

import (
	"github.com/goGo-service/back/internal/models"
	"github.com/goGo-service/back/internal/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user models.User) (int, error) {
	return s.repo.CreateUser(user)
}

func (s *UserService) GetUserByVkId(vkId int64) (*models.User, error) {
	return s.repo.GetUserByVkId(vkId)
}

func (s *UserService) GetUser(userId int) (*models.User, error) {
	user, err := s.repo.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}
