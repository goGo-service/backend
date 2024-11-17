package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/goGo-service/back/internal/models"
	"net/http"
	"strconv"
)

type createRoomRequest struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}

func (h *Handler) createRoom(c *gin.Context) {
	var requestBody createRoomRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}
	var newRoom models.Room
	newRoom.Name = requestBody.Name
	newRoom.Settings = models.RoomSettings{Capacity: 8}

	userID, exists := c.Get("UserId")
	if !exists {
		NewErrorResponse(c, http.StatusBadRequest, "user not found")
	}
	uId, ok := userID.(int)
	if !ok {
		NewErrorResponse(c, http.StatusBadRequest, "invalid user id")
	}
	id, err := h.roomUC.CreateNewRoom(newRoom, uId) //TODO: доставать из контекста
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(200, gin.H{
		"id":       id,
		"name":     newRoom.Name,
		"settings": newRoom.Settings,
	})
}

func (h *Handler) getRoom(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	userID, exists := c.Get("UserId")
	if !exists {
		NewErrorResponse(c, http.StatusBadRequest, "user not found")
	}
	uId, ok := userID.(int)
	if !ok {
		NewErrorResponse(c, http.StatusBadRequest, "invalid user id")
	}
	room, err := h.roomUC.GetRoom(id, uId) //TODO: доставать из контекста
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	//TODO: тут подумать, если нету прав или румы, по идее должно быть not access ошибка
	if room == nil {
		NewErrorResponse(c, http.StatusNotFound, "Room not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       id,
		"name":     room.Name,
		"settings": room.Settings,
	})
}
