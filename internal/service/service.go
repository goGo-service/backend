package service

import (
	goGO "github.com/goGo-service/back"
	"github.com/goGo-service/back/internal/repository"
)

type Authorization interface {
	CreateUser(user goGO.User) (int, error)
	GetUserByVkId(vkId int64) (*goGO.User, error)
}

type Service struct {
	Authorization
	VKAuth *AuthService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		VKAuth:        NewAuthService(repos.Authorization),
	}
}
