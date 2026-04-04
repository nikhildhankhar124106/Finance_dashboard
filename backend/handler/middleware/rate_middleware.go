package middleware

import (
	"log/slog"
	"net/http"
	"sync"
	"time"

	"backend/pkg/apperrors"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]*Client)
)

// cleanUpClients runs periodically to dynamically drop inactive map records freeing heap limits natively
func init() {
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
}

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	client, exists := clients[ip]
	if !exists {
		// 5 requests per second natively scaled limits, bursting to 10 maximum.
		limiter := rate.NewLimiter(5, 10)
		clients[ip] = &Client{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}

	client.lastSeen = time.Now()
	return client.limiter
}

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		limiter := getVisitor(clientIP)

		if !limiter.Allow() {
			slog.Warn("Rate Limit block hit tracking isolated limits", "ip", clientIP)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, apperrors.AppError{
				Status:  http.StatusTooManyRequests,
				Message: "Rate limit exceeded. Please block requests gracefully.",
			})
			return
		}

		c.Next()
	}
}
