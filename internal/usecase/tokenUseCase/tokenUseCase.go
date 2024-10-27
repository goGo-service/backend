package tokenUseCase

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"time"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

type TokenUseCase struct {
	redisClient *redis.Client
	secretKey   string
}

func NewTokenUseCase(redisClient *redis.Client) *TokenUseCase {
	return &TokenUseCase{
		redisClient: redisClient,
	}
}

func (u *TokenUseCase) RefreshToken(cookie string) (string, error) {
	if cookie != "refresh_token" {
		return "", fmt.Errorf("invalid refresh token")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		15,
	})
	ss, _ := token.SignedString([]byte(viper.GetString("SECRET_KEY")))

	return ss, nil
}
