package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"shutterspace/internal/domain"
)

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) BeginTx(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Begin()
}

// FindOverlapping mencari booking aktif yang overlap dengan rentang waktu yang diminta.
// Menggunakan PostgreSQL OVERLAPS operator dan SELECT FOR UPDATE NOWAIT
// untuk mencegah race condition pada concurrent requests.
//
// Dua rentang waktu (s1, e1) dan (s2, e2) overlap jika: s1 < e2 AND e1 > s2
func (r *bookingRepository) FindOverlapping(
	ctx context.Context,
	tx *gorm.DB,
	studioID uuid.UUID,
	date time.Time,
	startTime, endTime string,
) ([]domain.Booking, error) {
	var bookings []domain.Booking

	db := r.db
	if tx != nil {
		db = tx
	}

	err := db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE", Options: "NOWAIT"}).
		Where(`
			studio_id    = ?
			AND booking_date = ?::date
			AND status   IN ('pending', 'confirmed')
			AND (start_time, end_time) OVERLAPS (?::time, ?::time)
		`, studioID, date, startTime, endTime).
		Find(&bookings).Error

	return bookings, err
}

func (r *bookingRepository) FindByStudioAndDate(ctx context.Context, studioID uuid.UUID, date time.Time) ([]domain.Booking, error) {
	var bookings []domain.Booking
	err := r.db.WithContext(ctx).
		Where("studio_id = ? AND booking_date = ?::date AND status IN ('pending', 'confirmed')", studioID, date).
		Order("start_time ASC").
		Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) Create(ctx context.Context, tx *gorm.DB, booking *domain.Booking) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(booking).Error
}

func (r *bookingRepository) FindByUserID(ctx context.Context, userID uuid.UUID, status string, page, limit int) ([]domain.Booking, int64, error) {
	var bookings []domain.Booking
	var total int64

	query := r.db.WithContext(ctx).
		Preload("Studio").
		Preload("Studio.StudioType").
		Preload("Payment").
		Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Model(&domain.Booking{}).Count(&total)
	err := query.Order("created_at DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&bookings).Error

	return bookings, total, err
}

func (r *bookingRepository) FindByStudioID(ctx context.Context, studioID uuid.UUID, page, limit int) ([]domain.Booking, int64, error) {
	var bookings []domain.Booking
	var total int64

	query := r.db.WithContext(ctx).
		Preload("User").
		Preload("Payment").
		Where("studio_id = ?", studioID)

	query.Model(&domain.Booking{}).Count(&total)
	err := query.Order("booking_date DESC, start_time ASC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&bookings).Error

	return bookings, total, err
}

func (r *bookingRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	var booking domain.Booking
	err := r.db.WithContext(ctx).
		Preload("Studio").
		Preload("Studio.StudioType").
		Preload("User").
		Preload("Payment").
		First(&booking, "bookings.id = ?", id).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return &booking, nil
}

func (r *bookingRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.BookingStatus) error {
	return r.db.WithContext(ctx).
		Model(&domain.Booking{}).
		Where("id = ?", id).
		Update("status", status).Error
}
