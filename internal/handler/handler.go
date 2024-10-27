package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	goGO "github.com/goGo-service/back"
	"github.com/goGo-service/back/internal/models"
	"github.com/goGo-service/back/internal/service"

	"github.com/redis/go-redis/v9"
	"time"
)

type tokenUseCase interface {
	RefreshToken(oldToken string) (string, error)
}

type authUseCase interface {
	Auth(userId int) (*models.TokenPair, error)
	RefreshToken(oldToken string) (*models.TokenPair, error)
}

type ProfileUseCase interface {
	Profile(authHeader string) (*goGO.User, error)
}

type Handler struct {
	services      *service.Service
	RedisClient   *redis.Client
	vkAuthService *service.AuthService
	tokenUC       tokenUseCase
	authUC        authUseCase
	profileUC     ProfileUseCase
}

func NewHandler(services *service.Service, redisClient *redis.Client, vkAuthService *service.AuthService, tokenUC tokenUseCase, profileUC ProfileUseCase, authUC authUseCase) *Handler {
	return &Handler{
		services:      services,
		RedisClient:   redisClient,
		vkAuthService: vkAuthService,
		tokenUC:       tokenUC,
		profileUC:     profileUC,
		authUC:        authUC,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://welcome-satyr-easily.ngrok-free.app", "https://stallion-new-infinitely.ngrok-free.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Явно перечисляем методы
		AllowHeaders:     []string{"Content-Type", "Authorization"},           // Явно перечисляем заголовки
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           12 * time.Hour,
		AllowCredentials: true,
	}))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.GET("/logout", func(c *gin.Context) {
			c.JSON(200, "")
		})
		auth.GET("/redirect-url", h.redirectUrl)
		auth.GET("/token/refresh", h.refreshToken)
	}
	router.GET("/profile", h.profile)
	return router
}
