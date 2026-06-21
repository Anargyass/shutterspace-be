package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	pkgjwt "shutterspace/pkg/jwt"
	"shutterspace/pkg/response"
)

// RequireAuth memvalidasi JWT token dari Authorization header.
// Token yang valid akan meng-inject user_id, user_email, user_role ke dalam context.
func RequireAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Token autentikasi diperlukan")
			c.Abort()
			return
		}

		// Format: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Unauthorized(c, "Format token tidak valid. Gunakan: Bearer <token>")
			c.Abort()
			return
		}

		tokenStr := parts[1]

		claims, err := pkgjwt.VerifyToken(tokenStr, jwtSecret)
		if err != nil {
			response.Unauthorized(c, "Token tidak valid atau sudah kedaluwarsa")
			c.Abort()
			return
		}

		// Inject user info ke Gin context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}
