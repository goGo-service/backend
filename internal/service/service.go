package service

import (
	"github.com/goGo-service/back/internal/models"
	"github.com/goGo-service/back/internal/repository"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type Token interface {
	GenerateAccessToken(userId int, sessionID string) string
	GenerateRefreshToken(userId int, sessionID string) (string, error)
	RefreshTokens(refreshToken string) (*models.TokenPair, error)
	ParseToken(token string) (*models.TokenClaims, error)
	VerifyRefreshToken(refreshToken string, sessionID uuid.UUID) error
}

type User interface {
	CreateUser(user models.User) (int, error)
	GetUserByVkId(vkId int64) (*models.User, error)
	GetUser(userId int) (*models.User, error)
}

type VKID interface {
	GetUserInfo(accessToken string) (*models.VKIDUserInfo, error)
	ExchangeCode(code string, deviceId string, state string) (*VkidTokenResponse, error)
	CacheVKID(code string, id int64) error
	GetCachedVKID(code string) (int64, error)
	DeleteCachedVKID(code string) error
	GenerateStateAndCodeChallenge() (string, string, error)
}

type Service struct {
	User
	Token
	VKID
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:  NewUserService(repos.User),
		Token: NewTokenService(repos.User, viper.GetString("SECRET_KEY")),
		VKID:  NewVKIDService(repos.Cache),
	}
}
