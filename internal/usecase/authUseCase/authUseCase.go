package authUseCase

import (
	"github.com/goGo-service/back/internal/models"
	"github.com/goGo-service/back/internal/service"
	"github.com/google/uuid"
)

type AuthUseCase struct {
	services *service.Service
}

func NewAuthUseCase(service *service.Service) *AuthUseCase {
	return &AuthUseCase{
		services: service,
	}
}

func (u *AuthUseCase) Auth(userId int) (*models.TokenPair, error) {
	sessionID := uuid.New().String()
	accessToken := u.services.Token.GenerateAccessToken(userId, sessionID)
	refreshToken, err := u.services.Token.GenerateRefreshToken(userId, sessionID)
	if err != nil {
		return nil, err
	}
	tokenPair := &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return tokenPair, nil
}

func (u *AuthUseCase) RefreshToken(oldToken string) (*models.TokenPair, error) {
	//TODO: добавить проверку с бд
	tokens, err := u.services.Token.RefreshTokens(oldToken)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
