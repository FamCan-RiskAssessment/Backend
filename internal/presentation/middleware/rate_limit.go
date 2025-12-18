package middleware

import (
	"context"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/exception"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/ratelimit"
	"github.com/gin-gonic/gin"
)

type RateLimitMiddleware struct {
	constants   *bootstrap.Constants
	rateLimiter *ratelimit.RateLimiter
	security    *bootstrap.Security
}

func NewRateLimitMiddleware(
	constants *bootstrap.Constants,
	rateLimiter *ratelimit.RateLimiter,
	security *bootstrap.Security,
) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		constants:   constants,
		rateLimiter: rateLimiter,
		security:    security,
	}
}

func (rlm *RateLimitMiddleware) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		allowed, err := rlm.rateLimiter.Allow(context.Background(), ip)
		if err != nil {
			// Log error but don't block request on Redis failure
			// This is a fail-open approach
			c.Next()
			return
		}

		if !allowed {
			rateLimitErr := exception.NewRequestRateLimitError(
				"Too many requests from this IP",
				rlm.security.RateLimitPerMinute,
				nil,
			)
			panic(rateLimitErr) // Recovery middleware will handle this
		}

		c.Next()
	}
}
