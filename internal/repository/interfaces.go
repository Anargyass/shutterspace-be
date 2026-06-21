package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"shutterspace/internal/domain"
)

// ─── User Repository ───────────────────────────────────────────────────────

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

// ─── Studio Repository ─────────────────────────────────────────────────────

type StudioFilter struct {
	TypeSlug string
	Area     string
	MinPrice float64
	MaxPrice float64
	Page     int
	Limit    int
}

type StudioRepository interface {
	FindAll(ctx context.Context, filter StudioFilter) ([]domain.Studio, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Studio, error)
	FindBySlug(ctx context.Context, slug string) (*domain.Studio, error)
	Create(ctx context.Context, studio *domain.Studio) error
	Update(ctx context.Context, studio *domain.Studio) error
	FindByManagedBy(ctx context.Context, userID uuid.UUID) ([]domain.Studio, error)
}

// ─── Slot Repository ───────────────────────────────────────────────────────

type SlotRepository interface {
	FindByStudioAndDay(ctx context.Context, studioID uuid.UUID, day string) (*domain.AvailabilitySlot, error)
	FindAllByStudio(ctx context.Context, studioID uuid.UUID) ([]domain.AvailabilitySlot, error)
	UpsertSlots(ctx context.Context, slots []domain.AvailabilitySlot) error
}

// ─── Booking Repository ────────────────────────────────────────────────────

type BookingRepository interface {
	Create(ctx context.Context, tx *gorm.DB, booking *domain.Booking) error
	FindOverlapping(ctx context.Context, tx *gorm.DB, studioID uuid.UUID, date time.Time, startTime, endTime string) ([]domain.Booking, error)
	FindByStudioAndDate(ctx context.Context, studioID uuid.UUID, date time.Time) ([]domain.Booking, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, status string, page, limit int) ([]domain.Booking, int64, error)
	FindByStudioID(ctx context.Context, studioID uuid.UUID, page, limit int) ([]domain.Booking, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.BookingStatus) error
	BeginTx(ctx context.Context) *gorm.DB
}

// ─── Payment Repository ────────────────────────────────────────────────────

type PaymentRepository interface {
	Create(ctx context.Context, tx *gorm.DB, payment *domain.Payment) error
	FindByBookingID(ctx context.Context, bookingID uuid.UUID) (*domain.Payment, error)
	UpdateStatus(ctx context.Context, bookingID uuid.UUID, status domain.PaymentStatus, method domain.PaymentMethod, ref string) error
}
