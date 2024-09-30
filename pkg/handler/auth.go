package handler

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	goGO "github.com/goGo-service/back"
)

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
