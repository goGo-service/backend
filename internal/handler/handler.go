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
}

type Handler struct {
	services    *service.Service
	mw          Middleware
	redisClient *redis.Client
	authUC      authUseCase
	userUC      UserUseCase
	vkidUC      VKIDUseCase
	roomUC      RoomUseCase
}

func NewHandler(services *service.Service, mw Middleware, redisClient *redis.Client, userUC UserUseCase, authUC authUseCase, vkidUC VKIDUseCase, roomUC RoomUseCase) *Handler {
	return &Handler{
		services:    services,
		mw:          mw,
		redisClient: redisClient,
		userUC:      userUC,
		authUC:      authUC,
		vkidUC:      vkidUC,
		roomUC:      roomUC,
	}
}

// TODO добавить проверку на почту при регистрации
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(h.mw.CORS())

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.GET("/logout", h.mw.Auth(), h.logout)
		auth.GET("/redirect-url", h.redirectUrl)
		auth.GET("/token/refresh", h.refreshToken)
	}
	router.GET("/profile", h.mw.Auth(), h.profile)
	router.PATCH("/profile", h.mw.Auth(), h.editProfile)

	router.POST("/rooms", h.createRoom)

	return router
}
