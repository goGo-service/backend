package models

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// FIXME: как будто должно лежать в TokenService
type TokenClaims struct {
	jwt.RegisteredClaims
	UserId    int    `json:"user_id"`
	SessionId string `json:"session_id"`
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type RefreshToken struct {
	ID           int       `db:"id"`
	UserID       int       `db:"user_id"`
	SessionID    string    `db:"session_id"`
	RefreshToken string    `db:"refresh_token"`
	ExpireAt     time.Time `db:"expire_at"`
	CreatedAt    time.Time `db:"created_at"`
}
