package handler

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
	goGO "github.com/goGo-service/back"
)

type VKResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
}

//func (h *Handler) callbackHandler(c *gin.Context) {
//	code := c.Query("code")
//	if code == "" {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "No code in request"})
//		return
//	}
//
//	vkURL := fmt.Sprintf("%s?client_id=%s&client_secret=%s&redirect_uri=%s&code=%s",
//		tokenURL, clientID, clientSecret, redirectURI, code)
//
//	resp, err := http.Get(vkURL)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access token"})
//		return
//	}
//	defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
//		return
//	}
//
//	var vkResp VKResponse
//	if err := json.Unmarshal(body, &vkResp); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
//		return
//	}
//
//	userInfo, err := h.services.VKAuth.GetUserInfo(vkResp.AccessToken)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"access_token": vkResp.AccessToken,
//		"user_id":      vkResp.UserID,
//		"email":        vkResp.Email,
//		"user_info":    userInfo,
//	})
//}

//func (h *Handler) vkAuthHandler(c *gin.Context) {
//	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=%s&response_type=%s&v=%s",
//		vkAuthURL, clientID, redirectURI, scope, responseType, apiVersion)
//	c.Redirect(http.StatusFound, authURL)
//}

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
	//TODO: сохранять codeVerifier в кеш, где ключем будет state значением codeVerifier. +
	codeVerifier, _ := randomBytesInHex(32)
	//codeVerifier := "39365705206a4290cbf6b5aa1561ba8ab404b58df73ec30aceb823831dae38c7"
	sha2 := sha256.New()

	_, err := io.WriteString(sha2, codeVerifier)
	if err != nil {
		return
	}
	codeChallenge := base64.RawURLEncoding.EncodeToString(sha2.Sum(nil))
	state, _ := randomBytesInHex(24)
	//state := "9c00694677f5056d8060e6c43f847eda3bf08ba64a94827f"

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

func (h *Handler) signUp(c *gin.Context) {
	//TODO: проверять code из кеша, если все ок, то продолжаем. code:user_id

	var input goGO.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//FIXME: здесь не должен создаваться новый юзер, а изменяться старый
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

type signInRequestBody struct {
	Code     string `json:"code"`
	DeviceId string `json:"device_id"`
	State    string `json:"state"`
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

type vkidTokenResponse struct {
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

func (h *Handler) signIn(c *gin.Context) {
	//TODO: всю эту простыню кода привести в нормальный вид. -
	var requestBody signInRequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_text": "Invalid request"})
		return
	}
	//TODO: доставать state и codeVerifier из кеша и если че отдавать ошибку. +

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

	//TODO при невалидном токене code должен возвращать ошибку, а не пустой ответ). +
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
		newErrorResponse(c, http.StatusUnauthorized, "invalid request")
		newErrorResponse(c, http.StatusUnauthorized, "the provided request was invalid")
		return
	}

	fmt.Println(responseData)
	fmt.Println("ResponseData", responseData)
	user, err := h.services.GetUserByVkId(responseData.UserId)
	fmt.Println(responseData.AccessToken)
	if err != nil {
		//TODO при невалидном токене vkauth должен возвращать ошибку, а не пустого юзера). -
		userInfo, err := h.services.VKAuth.GetUserInfo(responseData.AccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error_text": "invalid access token"})
			return
		}
		if responseData.AccessToken != "" {
			newErrorResponse(c, http.StatusUnauthorized, "invalid request")
			newErrorResponse(c, http.StatusUnauthorized, "the provided request was invalid")
			return
		}

		user := goGO.User{
			FirstName: userInfo.Response[0].FirstName,
			VkID:      responseData.UserId,
			LastName:  userInfo.Response[0].LastName,
			Username:  "",
			Email:     userInfo.Response[0].Email,
		}
		_, err = h.services.CreateUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error_text": "failed to create user"})
			return
		}
		c.JSON(200, gin.H{
			"action":     "register",
			"first_name": userInfo.Response[0].FirstName,
			"last_name":  userInfo.Response[0].LastName,
			"email":      userInfo.Response[0].Email,
		})
		return
	}
	//TODO: сохранять code в кеш, чтобы при регистрации, можно было метчить.  +
	err = h.RedisClient.Set(context.Background(), requestBody.Code, responseData.UserId, 10&time.Minute).Err()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to save code in cache"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		user.Id,
	})
	ss, _ := token.SignedString([]byte("AllYourBase"))
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "tessts",
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Domain:   "localhost",
	})
	c.JSON(200, gin.H{
		"action":       "auth",
		"access_token": ss,
	})
}
