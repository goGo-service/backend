package userUseCase

import (
	goGO "github.com/goGo-service/back"
	"github.com/goGo-service/back/internal"
	"github.com/goGo-service/back/internal/service"
	"strings"
)

type UserUseCase struct {
	services *service.Service
}

func NewUserUseCase(service *service.Service) *UserUseCase {
	return &UserUseCase{
		services: service,
	}
}

func (u *UserUseCase) Profile(authHeader string) (*goGO.User, error) {
	accessToken := strings.Split(authHeader, " ")[1]

	if accessToken == "" {
		return nil, internal.AccessTokenRequiredError
	}

	user, err := u.services.GetUser(accessToken)
	if err != nil {
		return nil, internal.InternalServiceError
	}
	if user == nil {
		return nil, internal.UserNotFoundError
	}

	return user, nil
}
