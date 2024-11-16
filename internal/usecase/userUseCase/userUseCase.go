package userUseCase

import (
	"github.com/goGo-service/back/internal"
	"github.com/goGo-service/back/internal/models"
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

func (u *UserUseCase) GetByAccessToken(authHeader string) (*models.User, error) {
	accessToken := strings.Split(authHeader, " ")[1]

	if accessToken == "" {
		return nil, internal.AccessTokenRequiredError
	}
	tokenClaims, err := u.services.Token.ParseToken(accessToken)
	if err != nil {
		return nil, internal.AccessTokenRequiredError
	}
	user, err := u.services.User.GetUser(tokenClaims.UserId)
	if err != nil {
		return nil, internal.InternalServiceError
	}
	if user == nil {
		return nil, internal.UserNotFoundError
	}

	return user, nil
}

func (u *UserUseCase) CreateUser(user models.User) (int, error) {
	userId, err := u.services.User.CreateUser(user)
	if err != nil {
		return 0, err
	}

	return userId, err
}

func (u *UserUseCase) GetUserByVkId(id int64) (*models.User, error) {
	user, err := u.services.GetUserByVkId(id)
	if user == nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUseCase) UpdateUserFields(user *models.User, updates service.MutableUserFields) (bool, error) {

	return u.services.User.UpdateUserFields(user, updates)
}
