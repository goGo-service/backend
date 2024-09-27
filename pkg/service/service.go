package service

import (
	goGO "github.com/goGo-service/back"
	"github.com/goGo-service/back/pkg/repository"
)

type Authorization interface {
	CreateUser(user goGO.User) (int, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
