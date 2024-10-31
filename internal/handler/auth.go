package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/goGo-service/back/internal/models"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) redirectUrl(c *gin.Context) {
	res, err := h.vkidUC.GetRedirectUrl()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to save codeVerifier in cache"})
		return
	}

	scope := "email"
	c.JSON(200, gin.H{
		"app_id":         res.AppId,
		"redirect_url":   res.RedirectUrl,
		"state":          res.State,
		"code_challenge": res.CodeChallenge,
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
	id, vkidAT, err := h.vkidUC.GetUserIdAndAT(requestBody.Code, requestBody.DeviceId, requestBody.State)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "the provided request was invalid")
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
			c.JSON(http.StatusBadRequest, gin.H{"error_text": "invalid access token"})
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

func (h *Handler) refreshToken(c *gin.Context) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "token is required")
		return
	}

	tokens, err := h.authUC.RefreshToken(cookie)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "token expired or invalid")
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
