package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleUser        UserRole = "user"
	RoleStudioAdmin UserRole = "studio_admin"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name         string    `gorm:"type:varchar(100);not null"                     json:"name"`
	Email        string    `gorm:"type:varchar(150);not null;uniqueIndex"          json:"email"`
	PasswordHash string    `gorm:"type:text;not null"                             json:"-"`
	Phone        string    `gorm:"type:varchar(20)"                               json:"phone"`
	AvatarURL    string    `gorm:"type:text"                                      json:"avatar_url"`
	Role         UserRole  `gorm:"type:user_role;default:'user'"                  json:"role"`
	IsActive     bool      `gorm:"default:true"                                   json:"is_active"`
	CreatedAt    time.Time `                                                      json:"created_at"`
	UpdatedAt    time.Time `                                                      json:"updated_at"`
}
