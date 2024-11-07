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

type MutableUserFields struct {
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Username  string `db:"username" json:"username"`
}

func (s *UserService) UpdateUserFields(user *models.User, updates MutableUserFields) (bool, error) {
	isUpdated := false

	if updates.FirstName != "" && updates.FirstName != user.FirstName {
		user.FirstName = updates.FirstName
		isUpdated = true
	}
	if updates.LastName != "" && updates.LastName != user.LastName {
		user.LastName = updates.LastName
		isUpdated = true
	}
	if updates.Username != "" && updates.Username != user.Username {
		user.Username = updates.Username
		isUpdated = true
	}

	if isUpdated {
		return isUpdated, s.repo.UpdateUser(user)
	}

	return isUpdated, nil
}
