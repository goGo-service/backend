package service

import (
	"fmt"
	"github.com/goGo-service/back/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"time"
)

type TokenService struct {
	secretKey string
}

func NewTokenService(secretKey string) *TokenService {
	return &TokenService{secretKey: secretKey}
}

func (s *TokenService) GenerateAccessToken(userId int, sessionID string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		SessionId: sessionID,
		UserId:    userId,
	})
	ss, _ := token.SignedString([]byte(s.secretKey))
	return ss
}

func (s *TokenService) GenerateRefreshToken(userId int, sessionID string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		SessionId: sessionID,
		UserId:    userId,
	})
	ss, _ := token.SignedString([]byte(s.secretKey))
	return ss
}

func (s *TokenService) RefreshTokens(refreshToken string) (*models.TokenPair, error) {
	token, err := s.ParseToken(refreshToken)
	if err != nil {
		return nil, err
	}
	accessToken := s.GenerateAccessToken(token.UserId, token.SessionId)
	newRefreshToken := s.GenerateRefreshToken(token.UserId, token.SessionId)
	tokenPair := &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}
	return tokenPair, nil
}

func (s *TokenService) ParseToken(tokenString string) (*models.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		logrus.Print(s.secretKey)
		return nil, err
	}

	if claims, ok := token.Claims.(*models.TokenClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}

}
