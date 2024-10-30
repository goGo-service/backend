package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/goGo-service/back/internal"
	"github.com/goGo-service/back/internal/models"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"time"

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

	err = h.redisClient.Set(context.Background(), state, codeVerifier, 10&time.Minute).Err()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to save codeVerifier in cache"})
		return
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
	vkId, err := h.vkidUC.GetVKID(requestBody.Code)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid user_id")
		return
	}

	user, err := h.userUC.GetUserByVkId(vkId) // проверим, вдруг такой юзер уже есть
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if user != nil {
		newErrorResponse(c, http.StatusConflict, "user already exist")
		return
	}

	var input models.User
	input.Username = requestBody.Username
	input.FirstName = requestBody.FirstName
	input.LastName = requestBody.LastName
	input.VkID = vkId
	id, err := h.userUC.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	err = h.vkidUC.DeleteVKID(requestBody.Code)
	if err != nil {
		logrus.Warning("dont delete cached vkid", err.Error())
	}

	tokenPair, err := h.authUC.Auth(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	//TODO: вынести все генерации ответов с токеном в одну функцию
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Domain:   "localhost",
	})
	c.JSON(200, gin.H{
		"action":       "auth",
		"access_token": tokenPair.AccessToken,
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
	id, vkidAT, err := h.vkidUC.GetUserIdAndAT(requestBody.Code, requestBody.State, requestBody.DeviceId)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "the provided request was invalid")
		return
	}

	user, err := h.userUC.GetUserByVkId(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if user == nil {
		userInfo, err := h.vkidUC.GetUserInfo(vkidAT, requestBody.Code)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error_text": "invalid access token"})
			return
		}

		c.JSON(200, gin.H{
			"action":     "register",
			"first_name": userInfo.FirstName,
			"last_name":  userInfo.LastName,
			"email":      userInfo.Email,
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
	user, err := h.userUC.GetByAccessToken(authHeader)
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
