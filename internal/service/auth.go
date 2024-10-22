package service

import (
	"errors"
	"fmt"
	goGO "github.com/goGo-service/back"
	"github.com/goGo-service/back/internal/repository"
	"github.com/goccy/go-json"
	"github.com/golang-jwt/jwt/v5"
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
	return &AuthService{repo: repo}
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

func (s *AuthService) GetUser(accessToken string) (*goGO.User, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	var userId int64
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIdFloat, ok := claims["userId"].(float64)
		if !ok {
			return nil, errors.New("invalid token: userId not found")
		}
		userId = int64(userIdFloat)
	} else {
		return nil, errors.New("invalid token")
	}

	user, err := s.repo.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}
