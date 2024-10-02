package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/goGo-service/back/pkg/service"
	"time"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://stallion-new-infinitely.ngrok-free.app"},
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
	}

	router.GET("/profile", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"id":         1337,
			"username":   "screxy",
			"last_name":  "Миронов",
			"first_name": "Владислав",
			"email":      "dvbvladis@mail.ru",
		})
	})

	router.GET("/login", h.vkAuthHandler)
	router.GET("/callback", h.callbackHandler)

	return router
}
