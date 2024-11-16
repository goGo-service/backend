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
	GetUserById(c *gin.Context) (*models.User, error)
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

type Handler struct {
	services    *service.Service
	redisClient *redis.Client
	authUC      authUseCase
	userUC      UserUseCase
	vkidUC      VKIDUseCase
	mw          Middleware
}

func NewHandler(services *service.Service, redisClient *redis.Client, userUC UserUseCase, authUC authUseCase, vkidUC VKIDUseCase, mw Middleware) *Handler {
	return &Handler{
		services:    services,
		redisClient: redisClient,
		userUC:      userUC,
		authUC:      authUC,
		vkidUC:      vkidUC,
		mw:          mw,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(h.mw.CORS())

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.GET("/logout", h.mw.Auth(), h.logout)
		auth.GET("/redirect-url", h.redirectUrl)
		auth.GET("/token/refresh", h.mw.Auth(), h.refreshToken)
	}
	router.GET("/profile", h.mw.Auth(), h.profile)
	router.PATCH("/profile", h.mw.Auth(), h.editProfile)

	return router
}
