package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/goGo-service/back/internal/models"
	"github.com/goGo-service/back/internal/service"
	"github.com/redis/go-redis/v9"
	"time"
)

type authUseCase interface {
	Auth(userId int) (*models.TokenPair, error)
	RefreshToken(oldToken string) (*models.TokenPair, error)
}

type UserUseCase interface {
	GetByAccessToken(authHeader string) (*models.User, error)
	CreateUser(user models.User) (int, error)
	GetUserByVkId(id int64) (*models.User, error)
	UpdateUserFields(user *models.User, updates service.MutableUserFields) (bool, error)
}

type VKIDUseCase interface {
	GetUserIdAndAT(code string, deviceId string, state string) (int64, string, error)
	GetUserInfo(accessToken string, code string) (*models.VKIDUserInfo, error)
	GetVKID(code string) (int64, string, error)
	DeleteVKID(code string) error
	GetRedirectUrl() (*models.RedirectUrl, error)
}

type Handler struct {
	services    *service.Service
	redisClient *redis.Client
	authUC      authUseCase
	userUC      UserUseCase
	vkidUC      VKIDUseCase
}

func NewHandler(services *service.Service, redisClient *redis.Client, userUC UserUseCase, authUC authUseCase, vkidUC VKIDUseCase) *Handler {
	return &Handler{
		services:    services,
		redisClient: redisClient,
		userUC:      userUC,
		authUC:      authUC,
		vkidUC:      vkidUC,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://welcome-satyr-easily.ngrok-free.app", "https://stallion-new-infinitely.ngrok-free.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, // Явно перечисляем методы
		AllowHeaders:     []string{"Content-Type", "Authorization"},                    // Явно перечисляем заголовки
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           12 * time.Hour,
		AllowCredentials: true,
	}))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.GET("/logout", h.logout)
		auth.GET("/redirect-url", h.redirectUrl)
		auth.GET("/token/refresh", h.refreshToken)
	}
	router.GET("/profile", h.profile)
	router.PATCH("/profile", h.editProfile)

	return router
}
