package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/goGo-service/back/internal/service"
	"net/http"
)

type AuthMiddleware struct {
	services *service.Service
}

func NewAuthMiddleware(service *service.Service) *AuthMiddleware {
	return &AuthMiddleware{
		services: service,
	}
}

func (m *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		tokenClaims, err := m.services.Token.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("tokenClaims", tokenClaims)
		c.Next()
	}
}
