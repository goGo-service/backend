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
	UpdateUserFields(user *models.User, updates MutableUserFields) (bool, error)
}

type VKID interface {
	GetUserInfo(accessToken string) (*models.VKIDUserInfo, error)
	ExchangeCode(code string, deviceId string, state string) (*VkidTokenResponse, error)
	CacheVKID(code string, id int64, email string) error
	GetCachedVKID(code string) (int64, string, error)
	DeleteCachedVKID(code string) error
	GenerateStateAndCodeChallenge() (string, string, error)
}

type Room interface {
	CreateRoom(room models.Room) (int, error)
	AddOwnerToRoom(userId int, roomId int) error
	GetRoom(id int) (*models.Room, error)
	HaveAccess(userId int, roomId int) (bool, error)
	GetUserRooms(userID int) ([]*models.Room, error)
}

type Service struct {
	User
	Token
	VKID
	Room
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:  NewUserService(repos.User),
		Token: NewTokenService(viper.GetString("SECRET_KEY")),
		VKID:  NewVKIDService(repos.Cache),
		Room:  NewRoomService(repos.Room),
	}
}
