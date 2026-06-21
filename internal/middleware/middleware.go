package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ipLimiter menyimpan rate limiter per IP address
type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	limiters = make(map[string]*ipLimiter)
	mu       sync.Mutex
)

// RateLimiter membatasi request per IP address.
// r = request per detik, b = burst size
func RateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	// Goroutine untuk membersihkan limiter yang sudah tidak aktif
	go cleanupLimiters()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		if _, exists := limiters[ip]; !exists {
			limiters[ip] = &ipLimiter{
				limiter: rate.NewLimiter(r, b),
			}
		}
		limiters[ip].lastSeen = time.Now()
		limiter := limiters[ip].limiter
		mu.Unlock()

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "RATE_LIMIT_EXCEEDED",
				"message": "Terlalu banyak request. Coba lagi beberapa saat.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// cleanupLimiters menghapus entry limiter yang sudah tidak aktif > 3 menit
func cleanupLimiters() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, l := range limiters {
			if time.Since(l.lastSeen) > 3*time.Minute {
				delete(limiters, ip)
			}
		}
		mu.Unlock()
	}
}

// CORSMiddleware menambahkan header CORS yang aman.
// Hanya mengizinkan origin dari ALLOWED_ORIGIN env var.
func CORSMiddleware(allowedOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Hanya izinkan origin yang terdaftar — tidak ada wildcard *
		if origin == allowedOrigin {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			c.Header("Access-Control-Max-Age", "3600")
		}

		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
