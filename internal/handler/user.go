package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/goGo-service/back/internal"
	"net/http"
)

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

func (h *Handler) editProfile(c *gin.Context) {
	//TODO: ручка для изменения полей юзера
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
