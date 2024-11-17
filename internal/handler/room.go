package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/goGo-service/back/internal/models"
	"net/http"
)

type createRoomRequest struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}

func (h *Handler) createRoom(c *gin.Context) {
	var requestBody createRoomRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}
	var newRoom models.Room
	newRoom.Name = requestBody.Name
	newRoom.Settings = models.RoomSettings{Capacity: 8}
	id, err := h.roomUC.CreateNewRoom(newRoom, 19)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(200, gin.H{
		"id":      id,
		"name":    newRoom.Name,
		"setting": newRoom.Settings,
	})
}
