package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/goGo-service/back/internal/middleware/auth"
	"github.com/goGo-service/back/internal/middleware/security"
)

type Middleware interface {
	Auth() gin.HandlerFunc
	CORS() gin.HandlerFunc
}

type MwManager struct {
	authMiddleware *auth.AuthMiddleware
	corsMiddleware *security.CorsMiddleware
}

func NewMiddlewareManager(auth *auth.AuthMiddleware, cors *security.CorsMiddleware) *MwManager {
	return &MwManager{
		authMiddleware: auth,
		corsMiddleware: cors,
	}
}

func (m *MwManager) Auth() gin.HandlerFunc {
	return m.authMiddleware.Auth()
}

func (m *MwManager) CORS() gin.HandlerFunc { return m.corsMiddleware.CORS() }
