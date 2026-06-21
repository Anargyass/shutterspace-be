package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
)

type Booking struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null"                             json:"user_id"`
	StudioID       uuid.UUID      `gorm:"type:uuid;not null"                             json:"studio_id"`
	BookingDate    time.Time      `gorm:"type:date;not null"                             json:"booking_date"`
	StartTime      string         `gorm:"type:time;not null"                             json:"start_time"`
	EndTime        string         `gorm:"type:time;not null"                             json:"end_time"`
	DurationHours  float64        `gorm:"type:numeric(4,2);not null"                     json:"duration_hours"`
	Status         BookingStatus  `gorm:"type:booking_status;default:'pending'"          json:"status"`
	PricePerHour   float64        `gorm:"type:numeric(12,2);not null"                    json:"price_per_hour"`
	AddonsCost     float64        `gorm:"type:numeric(12,2);default:0"                   json:"addons_cost"`
	ServiceFee     float64        `gorm:"type:numeric(12,2);default:0"                   json:"service_fee"`
	TotalAmount    float64        `gorm:"type:numeric(12,2);not null"                    json:"total_amount"`
	SelectedAddons datatypes.JSON `gorm:"type:jsonb;default:'[]'"                        json:"selected_addons"`
	Notes          string         `gorm:"type:text"                                      json:"notes"`
	CancelledAt    *time.Time     `                                                      json:"cancelled_at"`
	CancelReason   string         `gorm:"type:text"                                      json:"cancel_reason"`
	CreatedAt      time.Time      `                                                      json:"created_at"`
	UpdatedAt      time.Time      `                                                      json:"updated_at"`

	// Relations
	User    *User    `gorm:"foreignKey:UserID"    json:"user,omitempty"`
	Studio  *Studio  `gorm:"foreignKey:StudioID"  json:"studio,omitempty"`
	Payment *Payment `gorm:"foreignKey:BookingID" json:"payment,omitempty"`
}
