package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/goGo-service/back/internal"
	"github.com/goGo-service/back/internal/service"
	"net/http"
)

func (h *Handler) profile(c *gin.Context) {
	user, err := h.userUC.GetUserById(c)
	if err != nil {
		switch {
		case errors.Is(err, internal.AccessTokenRequiredError):
			c.JSON(http.StatusBadRequest, gin.H{"error": "access token is required"})
		case errors.Is(err, internal.InternalServiceError):
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		case errors.Is(err, internal.UserNotFoundError):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		}
		return
	}
	//TODO: подумать норм ли это или нет
	c.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"username":   user.Username,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}

func (h *Handler) editProfile(c *gin.Context) {
	//TODO: ручка для изменения полей юзера
	user, err := h.userUC.GetUserById(c)
	if err != nil {
		switch {
		case errors.Is(err, internal.AccessTokenRequiredError):
			c.JSON(http.StatusBadRequest, gin.H{"error": "access token is required"})
		case errors.Is(err, internal.InternalServiceError):
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		case errors.Is(err, internal.UserNotFoundError):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		}
		return
	}
	var requestBody service.MutableUserFields
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.userUC.UpdateUserFields(user, requestBody)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
