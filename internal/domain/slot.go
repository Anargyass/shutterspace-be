package domain

import (
	"time"

	"github.com/google/uuid"
)

type DayOfWeek string

const (
	DayMonday    DayOfWeek = "monday"
	DayTuesday   DayOfWeek = "tuesday"
	DayWednesday DayOfWeek = "wednesday"
	DayThursday  DayOfWeek = "thursday"
	DayFriday    DayOfWeek = "friday"
	DaySaturday  DayOfWeek = "saturday"
	DaySunday    DayOfWeek = "sunday"
)

type AvailabilitySlot struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	StudioID   uuid.UUID `gorm:"type:uuid;not null"                             json:"studio_id"`
	DayOfWeek  DayOfWeek `gorm:"type:day_of_week;not null"                      json:"day"`
	OpenTime   string    `gorm:"type:time;default:'08:00:00'"                   json:"open"`
	CloseTime  string    `gorm:"type:time;default:'22:00:00'"                   json:"close"`
	IsOpen     bool      `gorm:"default:true"                                   json:"is_open"`
	CreatedAt  time.Time `                                                      json:"created_at"`
}
