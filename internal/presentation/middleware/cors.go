package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CORSMiddleware struct{}

func NewCorsMiddleware() *CORSMiddleware {
	return &CORSMiddleware{}
}

func (cm *CORSMiddleware) CORS() gin.HandlerFunc {
	corsConfig := cors.Config{
		// Remove AllowOrigins when using AllowOriginFunc
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "ngrok-skip-browser-warning"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowOriginFunc: func(origin string) bool {
			// Add some debugging to see what origin is being requested
			allowedOrigins := []string{"http://185.231.115.28:5173", "http://localhost:5173"}
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					return true
				}
			}
			// Log the origin that was rejected for debugging
			println("CORS: Rejected origin:", origin)
			return false
		},
	}

	return cors.New(corsConfig)
}
