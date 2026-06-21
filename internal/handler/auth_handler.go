package handler

import (
	"strings"

	"github.com/gin-gonic/gin"

	"shutterspace/internal/domain"
	"shutterspace/internal/service"
	"shutterspace/pkg/response"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(as service.AuthService) *AuthHandler {
	return &AuthHandler{authService: as}
}

type registerRequest struct {
	Name     string `json:"name"     binding:"required,min=2,max=100"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

type loginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		switch err {
		case domain.ErrEmailExists:
			response.Conflict(c, "EMAIL_EXISTS", "Email sudah terdaftar")
		default:
			response.InternalError(c, "Gagal melakukan registrasi")
		}
		return
	}

	response.Created(c, gin.H{"user": user})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}

	user, token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials, domain.ErrInactive:
			// SECURITY: Jangan beri informasi spesifik (email/password salah)
			response.Unauthorized(c, "Email atau password tidak valid")
		default:
			response.InternalError(c, "Gagal melakukan login")
		}
		return
	}

	response.OK(c, gin.H{
		"token": token,
		"user":  user,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Untuk JWT stateless, logout cukup dengan menghapus token di client.
	// Backend tidak menyimpan state token.
	// TODO(security): Implementasi token blacklist jika dibutuhkan invalidasi server-side.
	response.OK(c, gin.H{"message": "Logout berhasil"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userEmail, _ := c.Get("user_email")
	userRole, _ := c.Get("user_role")

	response.OK(c, gin.H{
		"user": gin.H{
			"id":    userID,
			"email": userEmail,
			"role":  userRole,
		},
	})
}

// ExtractBearerToken mengambil token dari Authorization header
func ExtractBearerToken(authHeader string) string {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}
	return ""
}
