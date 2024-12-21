package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/goGo-service/back/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"time"
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

func (h *Handler) getPresence(c *gin.Context) {
	userId, exists := c.Get("UserId")
	if !exists {
		NewErrorResponse(c, http.StatusBadRequest, "user not found")
		return
	}

	uid, ok := userId.(int)
	if !ok {
		NewErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}
	roomId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}
	presenceResponse, err := h.roomUC.GetRoomPresence(uid, roomId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "failed to fetch presence")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"room_presence": presenceResponse,
	})
}

func (h *Handler) getSubToken(c *gin.Context) {
	userId, exists := c.Get("UserId")
	if !exists {
		NewErrorResponse(c, http.StatusBadRequest, "user not found")
		return
	}

	uid, ok := userId.(int)
	if !ok {
		NewErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	roomId := c.Param("id")
	rId, err := strconv.Atoi(roomId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	haveRight, err := h.services.Room.HaveAccess(uid, rId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "error checking access")
		return
	}
	if !haveRight {
		NewErrorResponse(c, http.StatusForbidden, "you don`t have right")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.SubTokenClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3600 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Channel: "room:" + roomId,
		Sub:     strconv.Itoa(uid),
	})

	ss, err := token.SignedString([]byte(viper.GetString("CENTRIFUGO_TOKEN_HMAC_SECRET_KEY")))
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "failed to sign token")
		return
	}

	c.JSON(200, gin.H{
		"token": ss,
	})
}

type RoomPublishMessage struct {
	Type string `json:"type"`
}

func (h *Handler) roomMessage(c *gin.Context) {
	userId, exists := c.Get("UserId")
	if !exists {
		NewErrorResponse(c, http.StatusBadRequest, "user not found")
		return
	}

	uid, ok := userId.(int)
	if !ok {
		NewErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	roomId := c.Param("id")
	rId, err := strconv.Atoi(roomId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}
	var requestBody RoomPublishMessage
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}
	err = h.roomUC.PublishMessage(uid, rId, requestBody.Type)
	c.JSON(200, gin.H{})
}
