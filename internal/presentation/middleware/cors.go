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
		// TODO: Security Enhancement - Restrict CORS to specific origins
		// Currently allows all origins (*) which is a security risk in production
		// Action needed:
		// 1. Add ALLOWED_ORIGINS env variable in bootstrap/env.go
		// 2. Update this to use specific frontend domains:
		//    AllowOrigins: []string{"https://famcan.com", "https://app.famcan.com"}
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"POST", "GET", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization", "ngrok-skip-browser-warning"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}

	return cors.New(corsConfig)
}
