package service

import (
	goGO "github.com/goGo-service/back"
	"github.com/goGo-service/back/pkg/repository"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user goGO.User) (int, error) {
	return s.repo.CreateUser(user)
}
