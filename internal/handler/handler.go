package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/goGo-service/back/internal/models"
	"github.com/goGo-service/back/internal/service"
	"github.com/redis/go-redis/v9"
)

type authUseCase interface {
	Auth(userId int) (*models.TokenPair, error)
	RefreshToken(oldToken string) (*models.TokenPair, error)
}

type UserUseCase interface {
	GetUserById(id int) (*models.User, error)
	GetUserByVkId(id int64) (*models.User, error)
	CreateUser(user models.User) (int, error)
	UpdateUserFields(user *models.User, updates service.MutableUserFields) (bool, error)
}

type VKIDUseCase interface {
	GetUserIdAndAT(code string, deviceId string, state string) (int64, string, error)
	GetUserInfo(accessToken string, code string) (*models.VKIDUserInfo, error)
	GetVKID(code string) (int64, string, error)
	GetRedirectUrl() (*models.RedirectUrl, error)
	DeleteVKID(code string) error
}

type Middleware interface {
	Auth() gin.HandlerFunc
	CORS() gin.HandlerFunc
}

type RoomUseCase interface {
	CreateNewRoom(room models.Room, userId int) (int, error)
	GetRoom(roomId int, userId int) (*models.Room, error)
	GetUserRooms(userId int) ([]*models.Room, error)
	GetRoomPresence(userId int, roomId int) ([]*models.RoomUserPresence, error)
	PublishMessage(userId int, roomId int, message string) error
}

type Handler struct {
	mw          Middleware
	services    *service.Service
	redisClient *redis.Client
	authUC      authUseCase
	userUC      UserUseCase
	vkidUC      VKIDUseCase
	roomUC      RoomUseCase
}

func NewHandler(services *service.Service, mw Middleware, redisClient *redis.Client, userUC UserUseCase, authUC authUseCase, vkidUC VKIDUseCase, roomUC RoomUseCase) *Handler {
	return &Handler{
		services:    services,
		redisClient: redisClient,
		mw:          mw,
		userUC:      userUC,
		authUC:      authUC,
		vkidUC:      vkidUC,
		roomUC:      roomUC,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(h.mw.CORS())
	router.GET("/api/callback", func(c *gin.Context) {
		c.JSON(200, gin.H{})
		return
	})
	router.POST("/api/callback", func(c *gin.Context) {
		c.JSON(200, gin.H{})
		return
	})
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.GET("/logout", h.mw.Auth(), h.logout)
		auth.GET("/redirect-url", h.redirectUrl)
		auth.GET("/token/refresh", h.refreshToken)
		auth.GET("/token/connection", h.mw.Auth(), h.getConnToken)

	}
	users := router.Group("/api/v1/users", h.mw.Auth())
	{
		users.GET("/current", h.profile)
		users.PATCH("/current", h.editProfile)
		users.GET("/:id", h.getUser)
	}

	rooms := router.Group("/api/v1/rooms", h.mw.Auth())
	{
		rooms.POST("", h.createRoom)
		rooms.GET("", h.getUserRooms)
		rooms.GET("/:id", h.getRoom)
		rooms.POST("/:id/message", h.mw.Auth(), h.roomMessage)
		rooms.GET("/:id/presence", h.mw.Auth(), h.getPresence)
		rooms.GET("/:id/token/subscription", h.mw.Auth(), h.getSubToken)
	}

	return router
}
