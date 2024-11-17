package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/goGo-service/back/internal"
	"github.com/goGo-service/back/internal/service"
	"net/http"
)

func (h *Handler) profile(c *gin.Context) {
	userID, exists := c.Get("UserId")
	if !exists {
		NewErrorResponse(c, http.StatusUnauthorized, "user not found")
	}
	id, ok := userID.(int)
	if !ok {
		NewErrorResponse(c, http.StatusBadRequest, "invalid user id")
	}
	user, err := h.userUC.GetUserById(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.AccessTokenRequiredError):
			NewErrorResponse(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, internal.InternalServiceError):
			NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		case errors.Is(err, internal.UserNotFoundError):
			NewErrorResponse(c, http.StatusNotFound, err.Error())
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
	userID, exists := c.Get("UserId")
	if !exists {
		NewErrorResponse(c, http.StatusUnauthorized, "user not found")
	}
	id, ok := userID.(int)
	if !ok {
		NewErrorResponse(c, http.StatusBadRequest, "invalid user id")
	}
	user, err := h.userUC.GetUserById(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.AccessTokenRequiredError):
			NewErrorResponse(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, internal.InternalServiceError):
			NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		case errors.Is(err, internal.UserNotFoundError):
			NewErrorResponse(c, http.StatusNotFound, err.Error())
		}
		return
	}
	var requestBody service.MutableUserFields
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.userUC.UpdateUserFields(user, requestBody)

	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	NewErrorResponse(c, http.StatusOK, "User updated successfully")
}
