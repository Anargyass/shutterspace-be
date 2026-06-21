package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type SurabayaArea string

const (
	AreaPusat   SurabayaArea = "Surabaya Pusat"
	AreaBarat   SurabayaArea = "Surabaya Barat"
	AreaTimur   SurabayaArea = "Surabaya Timur"
	AreaSelatan SurabayaArea = "Surabaya Selatan"
	AreaUtara   SurabayaArea = "Surabaya Utara"
)

type StudioType struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(50);not null;uniqueIndex" json:"name"`
	Slug        string    `gorm:"type:varchar(50);not null;uniqueIndex" json:"slug"`
	IconURL     string    `gorm:"type:text" json:"icon_url"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Studio struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	StudioTypeID  int            `gorm:"not null"                                       json:"studio_type_id"`
	ManagedBy     *uuid.UUID     `gorm:"type:uuid"                                      json:"managed_by"`
	Name          string         `gorm:"type:varchar(150);not null"                     json:"name"`
	Slug          string         `gorm:"type:varchar(180);not null;uniqueIndex"          json:"slug"`
	Description   string         `gorm:"type:text"                                      json:"description"`
	Address       string         `gorm:"type:text;not null"                             json:"address"`
	Area          SurabayaArea   `gorm:"type:surabaya_area;not null"                    json:"area"`
	Capacity      int            `gorm:"default:1"                                      json:"capacity"`
	AreaSqm       float64        `gorm:"type:numeric(6,2)"                              json:"area_sqm"`
	PricePerHour  float64        `gorm:"type:numeric(12,2);not null"                    json:"price_per_hour"`
	Facilities    datatypes.JSON `gorm:"type:jsonb;default:'[]'"                        json:"facilities"`
	Addons        datatypes.JSON `gorm:"type:jsonb;default:'[]'"                        json:"addons"`
	Images        datatypes.JSON `gorm:"type:jsonb;default:'[]'"                        json:"images"`
	Rating        float64        `gorm:"type:numeric(3,2);default:0"                    json:"rating"`
	ReviewCount   int            `gorm:"default:0"                                      json:"review_count"`
	IsActive      bool           `gorm:"default:true"                                   json:"is_active"`
	CreatedAt     time.Time      `                                                      json:"created_at"`
	UpdatedAt     time.Time      `                                                      json:"updated_at"`

	// Relations (preload)
	StudioType        *StudioType        `gorm:"foreignKey:StudioTypeID"          json:"type,omitempty"`
	AvailabilitySlots []AvailabilitySlot `gorm:"foreignKey:StudioID"              json:"operating_hours,omitempty"`
}

// AddonItem digunakan untuk parse JSONB addons
type AddonItem struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
