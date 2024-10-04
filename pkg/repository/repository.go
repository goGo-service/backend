package repository

import (
	goGO "github.com/goGo-service/back"
	"github.com/goGo-service/back/pkg/repository/user"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user goGO.User) (int, error)
	GetUserByVkId(vkId int64) (*goGO.User, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: user.NewAuthPostgres(db),
	}
}
