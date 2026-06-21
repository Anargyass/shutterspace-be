package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{Success: true, Data: data})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{Success: true, Data: data})
}

func BadRequest(c *gin.Context, errCode, msg string) {
	c.JSON(http.StatusBadRequest, APIResponse{Success: false, Error: errCode, Message: msg})
}

func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, APIResponse{Success: false, Error: "UNAUTHORIZED", Message: msg})
}

func Forbidden(c *gin.Context, msg string) {
	c.JSON(http.StatusForbidden, APIResponse{Success: false, Error: "FORBIDDEN", Message: msg})
}

func NotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, APIResponse{Success: false, Error: "NOT_FOUND", Message: msg})
}

func Conflict(c *gin.Context, errCode, msg string) {
	c.JSON(http.StatusConflict, APIResponse{Success: false, Error: errCode, Message: msg})
}

func InternalError(c *gin.Context, msg string) {
	// Jangan expose detail error ke client
	c.JSON(http.StatusInternalServerError, APIResponse{Success: false, Error: "INTERNAL_ERROR", Message: msg})
}
