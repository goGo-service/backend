package service

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/goGo-service/back/internal/models"
	"github.com/goGo-service/back/internal/repository"
	"github.com/goccy/go-json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type VKIDService struct {
	appId           int
	vkidRedirectUrl string
	cache           repository.Cache
}

func NewVKIDService(cache repository.Cache) *VKIDService {
	return &VKIDService{appId: viper.GetInt("VKID_APP_ID"), vkidRedirectUrl: viper.GetString("VKID_REDIRECT_URL"), cache: cache}
}

type userInfoRequest struct {
	ClientId    int    `json:"client_id"`
	AccessToken string `json:"access_token"`
}

type userInfoResponse struct {
	models.VKIDUserInfo `json:"user"`
	Error               string `json:"error"`
	ErrorDescription    string `json:"error_description"`
	State               string `json:"state"`
}

func (s *VKIDService) GetUserInfo(accessToken string) (*models.VKIDUserInfo, error) {
	data := userInfoRequest{
		AccessToken: accessToken,
		ClientId:    s.appId,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		logrus.Fatalf("Ошибка сериализации данных: %v", err)
	}

	response, err := http.Post("https://id.vk.com/oauth2/user_info", "application/json", bytes.NewBuffer(jsonData))
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
	var responseData userInfoResponse

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		logrus.Fatalf("Ошибка парсинга JSON: %v", err)
	}

	if responseData.Error != "" || responseData.ErrorDescription != "" {
		return nil, fmt.Errorf("error: %s, desc: %s", responseData.Error, responseData.ErrorDescription)
	}

	return &responseData.VKIDUserInfo, nil
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

func (s *VKIDService) CacheVKID(code string, id int64, email string) error {
	vkidEmail := fmt.Sprintf("%d:%s", id, email)
	err := s.cache.Set(code, vkidEmail, 30*60)

	return err
}

func (s *VKIDService) GetCachedVKID(code string) (int64, string, error) {
	vkidEmail, err := s.cache.GetString(code)
	if err != nil {
		return 0, "", err
	}
	parts := strings.SplitN(vkidEmail, ":", 2)
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("invalid format: %s", vkidEmail)
	}

	vkid, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, "", fmt.Errorf("failed to parse vkid: %w", err)
	}

	email := parts[1]

	return vkid, email, nil
}

func (s *VKIDService) DeleteCachedVKID(code string) error {
	return s.cache.Delete(code)
}

func (s *VKIDService) GenerateStateAndCodeChallenge() (string, string, error) {
	codeVerifier, _ := randomBytesInHex(32)
	sha2 := sha256.New()

	_, err := io.WriteString(sha2, codeVerifier)
	if err != nil {
		return "", "", err
	}

	codeChallenge := base64.RawURLEncoding.EncodeToString(sha2.Sum(nil))
	state, err := randomBytesInHex(24)
	if err != nil {
		return "", "", err
	}
	err = s.cache.Set(state, codeVerifier, 30*60)
	if err != nil {
		return "", "", err
	}

	return state, codeChallenge, nil
}

func randomBytesInHex(count int) (string, error) {
	buf := make([]byte, count)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", fmt.Errorf("could not generate %d random bytes: %v", count, err)
	}

	return hex.EncodeToString(buf), nil
}
