package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"shutterspace/internal/domain"
	"shutterspace/internal/repository"
)

type PaymentService interface {
	ProcessMockPayment(ctx context.Context, bookingID, userID uuid.UUID, method domain.PaymentMethod) (*domain.Payment, error)
}

type paymentService struct {
	paymentRepo repository.PaymentRepository
	bookingRepo repository.BookingRepository
}

func NewPaymentService(paymentRepo repository.PaymentRepository, bookingRepo repository.BookingRepository) PaymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
		bookingRepo: bookingRepo,
	}
}

// ProcessMockPayment mensimulasikan pembayaran sukses (untuk demo/prototipe).
// Dalam implementasi nyata, ini akan memanggil Midtrans/Xendit SDK.
func (s *paymentService) ProcessMockPayment(ctx context.Context, bookingID, userID uuid.UUID, method domain.PaymentMethod) (*domain.Payment, error) {
	booking, err := s.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	if booking.UserID != userID {
		return nil, domain.ErrForbidden
	}

	if booking.Status != domain.BookingStatusPending {
		return nil, fmt.Errorf("%w: booking sudah diproses atau dibatalkan", domain.ErrConflict)
	}

	// Generate referensi simulasi
	externalRef := fmt.Sprintf("MOCK-TXN-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)

	// Update payment status → paid
	if err := s.paymentRepo.UpdateStatus(ctx, bookingID, domain.PaymentStatusPaid, method, externalRef); err != nil {
		return nil, err
	}

	// Update booking status → confirmed
	if err := s.bookingRepo.UpdateStatus(ctx, bookingID, domain.BookingStatusConfirmed); err != nil {
		return nil, err
	}

	// Return updated payment
	payment, err := s.paymentRepo.FindByBookingID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	return payment, nil
}
