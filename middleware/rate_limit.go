package middleware

import (
	"net/http"
	"sync"
	"time"

	"AUV/config"

	"github.com/gin-gonic/gin"
)

type tokenBucket struct {
	IP          string
	LastRequest time.Time
	Tokens      int
}

var (
	buckets = make(map[string]*tokenBucket)
	mutex   sync.Mutex
)

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		if forwarded := c.GetHeader("X-Forwarded-For"); forwarded != "" {
			clientIP = forwarded
		}

		mutex.Lock()
		defer mutex.Unlock()

		bucket, exists := buckets[clientIP]
		if !exists {
			bucket = &tokenBucket{
				IP:          clientIP,
				Tokens:      config.Cfg.Server.LimitRate,
				LastRequest: time.Now(),
			}
			buckets[clientIP] = bucket
		} else {
			elapsed := time.Since(bucket.LastRequest)
			if elapsed >= time.Second {
				bucket.Tokens = config.Cfg.Server.LimitRate
				bucket.LastRequest = time.Now()
			}
		}

		if bucket.Tokens <= 0 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			return
		}

		bucket.Tokens--
		c.Next()
	}
}
