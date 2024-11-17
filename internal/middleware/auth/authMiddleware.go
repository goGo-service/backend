package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/goGo-service/back/internal/handler"
	"github.com/goGo-service/back/internal/service"
	"net/http"
	"strings"
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
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			handler.NewErrorResponse(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			handler.NewErrorResponse(c, http.StatusUnauthorized, "Token is missing")
			c.Abort()
			return
		}

		tokenClaims, err := m.services.Token.ParseToken(token)
		if err != nil {
			handler.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Set("UserId", tokenClaims.UserId)
		c.Set("SessionId", tokenClaims.SessionId)
		c.Next()
	}
}
