package service

import (
	"errors"
	"fmt"
	goGO "github.com/goGo-service/back"
	"github.com/goGo-service/back/internal/handler"
	"github.com/goGo-service/back/internal/repository"
	"github.com/goccy/go-json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"net/http"
)

const (
	apiURL     = "https://api.vk.com/method/users.get"
	apiVersion = "5.199"
)

type UserResponse struct {
	Response []struct {
		ID        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	} `json:"response"`
}

type AuthService struct {
	repo      repository.Authorization
	jwtSecret string
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo, jwtSecret: viper.GetString("SECRET_KEY")}
}

func (s *AuthService) CreateUser(user goGO.User) (int, error) {
	return s.repo.CreateUser(user)
}

func (s *AuthService) GetUserInfo(accessToken string) (*UserResponse, error) {
	vkAPIURL := fmt.Sprintf("%s?access_token=%s&v=%s", apiURL, accessToken, apiVersion)
	resp, err := http.Get(vkAPIURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userResp UserResponse
	if err := json.Unmarshal(body, &userResp); err != nil {
		return nil, err
	}

	return &userResp, nil
}

func (s *AuthService) GetUserByVkId(vkId int64) (*goGO.User, error) {
	return s.repo.GetUserByVkId(vkId)
}

func ParseToken(tokenString string, secretKey string) (*handler.TokenClaims, error) {
	// Разбираем и валидируем токен
	token, err := jwt.ParseWithClaims(tokenString, &handler.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Приводим claims к типу TokenClaims
	if claims, ok := token.Claims.(*handler.TokenClaims); ok && token.Valid {
		return claims, nil
	} else {
		logrus.Debug("her")
		return nil, errors.New("invalid token")
	}
}
func (s *AuthService) GetUser(accessToken string) (*goGO.User, error) {
	claims, err := ParseToken(accessToken, viper.GetString("SECRET_KEY"))
	if err != nil {

		// TODO: Обрабатываем ошибку
	}

	userId := claims.UserId

	user, err := s.repo.GetUserById(int64(userId))
	if err != nil {
		return nil, err
	}

	return user, nil
}
