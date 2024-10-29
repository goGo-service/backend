package handler

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/goGo-service/back/internal"
	"github.com/goGo-service/back/internal/models"
	"github.com/goccy/go-json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
)

type VKResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
}

func randomBytesInHex(count int) (string, error) {
	buf := make([]byte, count)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", fmt.Errorf("could not generate %d random bytes: %v", count, err)
	}

	return hex.EncodeToString(buf), nil
}

func (h *Handler) redirectUrl(c *gin.Context) {
	redirectUrl := viper.GetString("VKID_REDIRECT_URL")
	appId := viper.GetString("VKID_APP_ID")
	codeVerifier, _ := randomBytesInHex(32)
	sha2 := sha256.New()

	_, err := io.WriteString(sha2, codeVerifier)
	if err != nil {
		return
	}
	codeChallenge := base64.RawURLEncoding.EncodeToString(sha2.Sum(nil))
	state, _ := randomBytesInHex(24)

	err = h.RedisClient.Set(context.Background(), state, codeVerifier, 10&time.Minute).Err()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to save codeVerifier in cache"})
	}

	scope := "email"
	c.JSON(200, gin.H{
		"app_id":         appId,
		"redirect_url":   redirectUrl,
		"state":          state,
		"code_challenge": codeChallenge,
		"scope":          scope,
	})
}

type SignUpRequestBody struct {
	Code      string `json:"code" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Username  string `json:"username" binding:"required"`
}

func (h *Handler) signUp(c *gin.Context) {
	var requestBody SignUpRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	vkUserId, err := h.RedisClient.Get(context.Background(), requestBody.Code).Result()
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid user_id")
		return
	}

	var input models.User
	input.Username = requestBody.Username
	input.FirstName = requestBody.FirstName
	input.LastName = requestBody.LastName
	//FIXME: выглядит дерьмово
	vkId, _ := strconv.Atoi(vkUserId)
	input.VkID = int64(vkId)
	id, err := h.services.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	h.RedisClient.Del(context.Background(), requestBody.Code)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserId: id,
	})
	ss, _ := token.SignedString([]byte(viper.GetString("SECRET_KEY")))
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "refresh_token",
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Domain:   "localhost",
	})
	c.JSON(200, gin.H{
		"action":       "auth",
		"access_token": ss,
	})
}

type signInRequestBody struct {
	Code     string `json:"code"`
	DeviceId string `json:"device_id"`
	State    string `json:"state"`
}

func (h *Handler) signIn(c *gin.Context) {
	var requestBody signInRequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_text": "Invalid request"})
		return
	}

	codeVerifier, err := h.RedisClient.Get(context.Background(), requestBody.State).Result()
	if err != nil || codeVerifier == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid state")
	}

	data := vkidTokenRequest{
		ClientId:     viper.GetInt("VKID_APP_ID"),
		GrantType:    "authorization_code",
		DeviceId:     requestBody.DeviceId,
		Code:         requestBody.Code,
		CodeVerifier: codeVerifier,
		RedirectUri:  viper.GetString("VKID_REDIRECT_URL"),
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
	var responseData vkidTokenResponse

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		logrus.Fatalf("Ошибка парсинга JSON: %v", err)
	}

	if responseData.Error != "" || responseData.ErrorDescription != "" {
		newErrorResponse(c, http.StatusUnauthorized, "the provided request was invalid")
		return
	}

	user, err := h.services.GetUserByVkId(responseData.UserId)
	if err != nil {
		userInfo, err := h.services.VKID.GetUserInfo(responseData.AccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error_text": "invalid access token"})
			return
		}
		err = h.RedisClient.Set(context.Background(), requestBody.Code, responseData.UserId, 10&time.Minute).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		c.JSON(200, gin.H{
			"action":     "register",
			"first_name": userInfo.Response[0].FirstName,
			"last_name":  userInfo.Response[0].LastName,
			"email":      userInfo.Response[0].Email,
		})
		return
	}

	token, err := h.authUC.Auth(user.Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    token.RefreshToken,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Domain:   "localhost",
	})
	c.JSON(200, gin.H{
		"action":       "auth",
		"access_token": token.AccessToken,
	})
}

func (h *Handler) profile(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}
	user, err := h.profileUC.Profile(authHeader)
	if err != nil {
		switch err {
		case internal.AccessTokenRequiredError:
			c.JSON(http.StatusBadRequest, gin.H{"error": "access token is required"})
		case internal.InternalServiceError:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		case internal.UserNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) refreshToken(c *gin.Context) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		return
	}

	tokens, err := h.authUC.RefreshToken(cookie)
	if err != nil {
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Domain:   "localhost",
	})
	c.JSON(200, gin.H{
		"access_token": tokens.AccessToken,
	})
}
