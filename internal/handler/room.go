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
	uid, ok := userID.(int)
	if !ok {
		NewErrorResponse(c, http.StatusBadRequest, "invalid user id")
	}
	id, err := h.roomUC.CreateNewRoom(newRoom, uid) //TODO: доставать из контекста
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	newRoom.Id = id
	c.JSON(http.StatusOK, newRoom.ToResponse())
}

func (h *Handler) getRoom(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	userID, exists := c.Get("UserId")
	if !exists {
		NewErrorResponse(c, http.StatusBadRequest, "user not found")
		return
	}
	uid, ok := userID.(int)
	if !ok {
		NewErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}
	room, err := h.roomUC.GetRoom(id, uid) //TODO: доставать из контекста
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	//TODO: тут подумать, если нету прав или румы, по идее должно быть not access ошибка
	if room == nil {
		NewErrorResponse(c, http.StatusNotFound, "Room not found")
		return
	}

	c.JSON(http.StatusOK, room.ToResponse())
}

func (h *Handler) getUserRooms(c *gin.Context) {
	userID, exists := c.Get("UserId")
	if !exists {
		NewErrorResponse(c, http.StatusBadRequest, "user not found")
		return
	}
	uid, ok := userID.(int)
	if !ok {
		NewErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}
	rooms, err := h.roomUC.GetUserRooms(uid) //TODO: доставать из контекста
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}
	var responseRooms []models.RoomResponse
	for _, room := range rooms {
		responseRooms = append(responseRooms, *room.ToResponse()) // Вызов метода ToResponse
	}

	c.JSON(http.StatusOK, gin.H{
		"rooms": responseRooms,
	})
}
