package userUseCase

import (
	"github.com/gin-gonic/gin"
	"github.com/goGo-service/back/internal"
	"github.com/goGo-service/back/internal/models"
	"github.com/goGo-service/back/internal/service"
)

type UserUseCase struct {
	services *service.Service
}

func NewUserUseCase(service *service.Service) *UserUseCase {
	return &UserUseCase{
		services: service,
	}
}

func (u *UserUseCase) GetUserById(c *gin.Context) (*models.User, error) {
	userID, exists := c.Get("UserId")
	if !exists {
		return nil, internal.AccessTokenRequiredError
	}
	id, ok := userID.(int)
	if !ok {
		return nil, internal.InvalidUserIDError
	}
	user, err := u.services.User.GetUser(id)
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
