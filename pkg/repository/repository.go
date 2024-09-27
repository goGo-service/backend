package repository

import (
	goGO "github.com/goGo-service/back"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user goGO.User) (int, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
