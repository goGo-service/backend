package service

import (
	"bytes"
	"fmt"
	"github.com/goGo-service/back/internal/models"
	"github.com/goGo-service/back/internal/repository"
	"github.com/goccy/go-json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
)

const (
	apiURL     = "https://api.vk.com/method/users.get"
	apiVersion = "5.199"
)

type UserResponse struct {
	Response []struct {
		models.VKIDUserInfo
	} `json:"response"`
}

type VKIDService struct {
	appId           int
	vkidRedirectUrl string
	cache           repository.Cache
}

func NewVKIDService(cache repository.Cache) *VKIDService {
	return &VKIDService{appId: viper.GetInt("VKID_APP_ID"), vkidRedirectUrl: viper.GetString("VKID_REDIRECT_URL"), cache: cache}
}

func (s *VKIDService) GetUserInfo(accessToken string) (*UserResponse, error) {
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

type vkidTokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientId     int    `json:"client_id"`
	DeviceId     string `json:"device_id"`
	RedirectUri  string `json:"redirect_uri"`
	Code         string `json:"code"`
	CodeVerifier string `json:"code_verifier"`
	Scope        string `json:"scope"`
}

type VkidTokenResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	IdToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	UserId       int64  `json:"user_id"`
	State        string `json:"state"`
	Scope        string `json:"scope"`
	errors
}

type errors struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (s *VKIDService) ExchangeCode(code string, deviceId string, state string) (*VkidTokenResponse, error) {
	codeVerifier, err := s.cache.GetString(state)
	if err != nil {
		return nil, err
	}

	data := vkidTokenRequest{
		ClientId:     s.appId,
		GrantType:    "authorization_code",
		Code:         code,
		DeviceId:     deviceId,
		CodeVerifier: codeVerifier,
		RedirectUri:  s.vkidRedirectUrl,
		Scope:        "email",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		logrus.Fatalf("Ошибка сериализации данных: %v", err)
	}

	response, err := http.Post("https://id.vk.com/oauth2/auth", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.Fatalf("failed to send request: %v", err)
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Fatalf("Ошибка чтения ответа: %v", err)
	}
	var responseData VkidTokenResponse

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		logrus.Fatalf("Ошибка парсинга JSON: %v", err)
	}

	if responseData.Error != "" || responseData.ErrorDescription != "" {
		return nil, fmt.Errorf(responseData.Error, responseData.ErrorDescription)
	}
	return &responseData, nil
}

func (s *VKIDService) CacheVKIDUser(code string, id int64) error {
	err := s.cache.Set(code, id, 30*60)
	return err
}
