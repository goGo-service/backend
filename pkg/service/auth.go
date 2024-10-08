package service

import (
	"fmt"
	goGO "github.com/goGo-service/back"
	"github.com/goGo-service/back/pkg/repository"
	"github.com/goccy/go-json"
	"io/ioutil"
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
	repo repository.Authorization
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

	body, err := ioutil.ReadAll(resp.Body)
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
