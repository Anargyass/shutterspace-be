package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"shutterspace/internal/domain"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, tx *gorm.DB, payment *domain.Payment) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(payment).Error
}

func (r *paymentRepository) FindByBookingID(ctx context.Context, bookingID uuid.UUID) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.WithContext(ctx).
		Where("booking_id = ?", bookingID).
		First(&payment).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return &payment, nil
}

func (r *paymentRepository) UpdateStatus(
	ctx context.Context,
	bookingID uuid.UUID,
	status domain.PaymentStatus,
	method domain.PaymentMethod,
	ref string,
) error {
	updates := map[string]interface{}{
		"status":         status,
		"payment_method": method,
		"external_ref":   ref,
	}
	if status == domain.PaymentStatusPaid {
		now := time.Now()
		updates["paid_at"] = &now
	}

	result := r.db.WithContext(ctx).
		Model(&domain.Payment{}).
		Where("booking_id = ?", bookingID).
		Updates(updates)

	if result.RowsAffected == 0 {
		return fmt.Errorf("payment for booking %s not found", bookingID)
	}
	return result.Error
}
