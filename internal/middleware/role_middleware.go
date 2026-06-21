package middleware

import (
	"github.com/gin-gonic/gin"

	"shutterspace/pkg/response"
)

// RequireRole memastikan user yang sudah autentikasi memiliki role yang dibutuhkan.
// Harus digunakan SETELAH RequireAuth middleware.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			response.Unauthorized(c, "Autentikasi diperlukan")
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			response.Unauthorized(c, "Data role tidak valid")
			c.Abort()
			return
		}

		for _, allowed := range roles {
			if roleStr == allowed {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "Anda tidak memiliki izin untuk mengakses resource ini")
		c.Abort()
	}
}
