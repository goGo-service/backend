package handler

import (
	"fmt"
	"github.com/goccy/go-json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
	goGO "github.com/goGo-service/back"
)

// TODO заполнить данные после регистрации goGO в вк
const (
	clientID     = ""
	redirectURI  = ""
	vkAuthURL    = "https://oauth.vk.com/authorize"
	scope        = "email"
	responseType = "code"
	apiVersion   = "5.131"
	clientSecret = ""
	tokenURL     = "https://oauth.vk.com/access_token"
)

type VKResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
}

func (h *Handler) callbackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code in request"})
		return
	}

	vkURL := fmt.Sprintf("%s?client_id=%s&client_secret=%s&redirect_uri=%s&code=%s",
		tokenURL, clientID, clientSecret, redirectURI, code)

	resp, err := http.Get(vkURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access token"})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	var vkResp VKResponse
	if err := json.Unmarshal(body, &vkResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	userInfo, err := h.services.VKAuth.GetUserInfo(vkResp.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": vkResp.AccessToken,
		"user_id":      vkResp.UserID,
		"email":        vkResp.Email,
		"user_info":    userInfo,
	})
}

func (h *Handler) vkAuthHandler(c *gin.Context) {
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=%s&response_type=%s&v=%s",
		vkAuthURL, clientID, redirectURI, scope, responseType, apiVersion)
	c.Redirect(http.StatusFound, authURL)
}

func (h *Handler) signUp(c *gin.Context) {
	var input goGO.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

func (h *Handler) signIn(c *gin.Context) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		1,
	})
	ss, _ := token.SignedString([]byte("AllYourBase"))
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "tesst",
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: false,
		Domain:   "stallion-new-infinitely.ngrok-free.app",
		//SameSite: http.SameSiteNoneMode,
		//Secure:   true,
	})
	c.JSON(200, gin.H{
		"action":       "auth",
		"access_token": ss,
	})
}
