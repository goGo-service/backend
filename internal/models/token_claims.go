package models

import "github.com/golang-jwt/jwt/v5"

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
