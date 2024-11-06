package repository

import (
	"github.com/goGo-service/back/internal/models"
	"github.com/goGo-service/back/internal/repository/cache"
	"github.com/goGo-service/back/internal/repository/user"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type User interface {
	CreateUser(user models.User) (int, error)
	GetUserByVkId(vkId int64) (*models.User, error)
	GetUserById(userId int) (*models.User, error)
	SaveRefreshToken(token models.RefreshToken) error
}

type Cache interface {
	GetString(key string) (string, error)
	GetInt(key string) (int, error)
	GetInt64(key string) (int64, error)
	Set(key string, value any, ttl int) error
	Delete(key string) error
}

type Repository struct {
	User
	Cache
}

func NewRepository(db *sqlx.DB, cacheClient *redis.Client) *Repository {
	return &Repository{
		User:  user.NewUserPostgres(db),
		Cache: cache.NewRedisCache(cacheClient),
	}
}
