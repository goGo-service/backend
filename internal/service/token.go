package service

import (
	"fmt"
	"github.com/goGo-service/back/internal/models"
	"github.com/goGo-service/back/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type TokenService struct {
	secretKey string
	repo      repository.User
}

func NewTokenService(repo repository.User, secretKey string) *TokenService {
	return &TokenService{secretKey: secretKey, repo: repo}
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

func (s *TokenService) GenerateRefreshToken(userId int, sessionID string) (string, error) {
	expireAt := time.Now().Add(30 * 24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		SessionId: sessionID,
		UserId:    userId,
	})
	ss, _ := token.SignedString([]byte(s.secretKey))
	refreshToken := models.RefreshToken{
		UserID:       userId,
		SessionID:    sessionID,
		RefreshToken: ss,
		ExpireAt:     expireAt,
	}
	err := s.repo.SaveRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}
	return ss, nil
}

func (s *TokenService) RefreshTokens(refreshToken string) (*models.TokenPair, error) {
	token, err := s.ParseToken(refreshToken)
	if err != nil {
		return nil, err
	}
	accessToken := s.GenerateAccessToken(token.UserId, token.SessionId)
	newRefreshToken, err := s.GenerateRefreshToken(token.UserId, token.SessionId)
	if err != nil {
		return nil, err
	}
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
		return nil, err
	}

	if claims, ok := token.Claims.(*models.TokenClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}

}

func (s *TokenService) VerifyRefreshToken(refreshToken string, sessionID uuid.UUID) error {
	tokenData, err := s.repo.GetRefreshToken(refreshToken, sessionID)
	if err != nil {
		return err
	}

	if time.Now().After(tokenData.ExpireAt) {
		fmt.Println("Token expired at:", tokenData.ExpireAt)
		return fmt.Errorf("invalid or expired token")
	}

	return nil
}
