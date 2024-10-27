package service

import (
	goGO "github.com/goGo-service/back"
	"github.com/goGo-service/back/internal/models"
	"github.com/goGo-service/back/internal/repository"
	"github.com/spf13/viper"
)

type Authorization interface {
	CreateUser(user goGO.User) (int, error)
	GetUserByVkId(vkId int64) (*goGO.User, error)
	GetUser(accessToken string) (*goGO.User, error)
}

type Token interface {
	GenerateAccessToken(userId int, sessionID string) string
	GenerateRefreshToken(userId int, sessionID string) string
	RefreshTokens(refreshToken string) (*models.TokenPair, error)
	ParseToken(token string) (*models.TokenClaims, error)
}

type Service struct {
	Authorization
	VKAuth *AuthService
	Token  *TokenService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		VKAuth:        NewAuthService(repos.Authorization),
		Token:         NewTokenService(viper.GetString("SECRET_KEY")),
	}
}
