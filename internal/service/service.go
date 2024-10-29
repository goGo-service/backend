package service

import (
	"github.com/goGo-service/back/internal/models"
	"github.com/goGo-service/back/internal/repository"
	"github.com/spf13/viper"
)

type Token interface {
	GenerateAccessToken(userId int, sessionID string) string
	GenerateRefreshToken(userId int, sessionID string) string
	RefreshTokens(refreshToken string) (*models.TokenPair, error)
	ParseToken(token string) (*models.TokenClaims, error)
}

type User interface {
	CreateUser(user models.User) (int, error)
	GetUserByVkId(vkId int64) (*models.User, error)
	GetUser(userId int) (*models.User, error)
}

type VKID interface {
	GetUserInfo(accessToken string) (*UserResponse, error)
	ExchangeCode(code string, deviceId string, state string) (*VkidTokenResponse, error)
}

type Service struct {
	User
	Token
	VKID
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:  NewUserService(repos.User),
		Token: NewTokenService(viper.GetString("SECRET_KEY")),
		VKID:  NewVKIDService(repos.Cache),
	}
}
