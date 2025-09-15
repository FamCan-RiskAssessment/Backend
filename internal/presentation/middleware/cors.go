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
		AllowOriginFunc: func(origin string) bool {
			// Allow localhost for development
			if origin == "http://localhost:5173" || origin == "http://localhost:3000" {
				return true
			}
			// Allow your production origin
			if origin == "http://185.231.115.28:5173" {
				return true
			}
			return false
		},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "ngrok-skip-browser-warning"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}

	return cors.New(corsConfig)
}
