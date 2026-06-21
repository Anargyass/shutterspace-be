package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"shutterspace/internal/domain"
	"shutterspace/internal/repository"
)

type CreateBookingInput struct {
	UserID         uuid.UUID
	StudioID       uuid.UUID
	BookingDate    time.Time
	StartTime      string  // "11:00"
	DurationHours  float64 // minimum 1.0
	SelectedAddons []domain.AddonItem
	Notes          string
}

type BookingService interface {
	CreateBooking(ctx context.Context, input CreateBookingInput) (*domain.Booking, error)
	GetUserBookings(ctx context.Context, userID uuid.UUID, status string, page, limit int) ([]domain.Booking, int64, error)
	GetBookingByID(ctx context.Context, bookingID, requesterID uuid.UUID, requesterRole string) (*domain.Booking, error)
	CancelBooking(ctx context.Context, bookingID, userID uuid.UUID) error
	GetStudioBookings(ctx context.Context, studioID, adminID uuid.UUID, page, limit int) ([]domain.Booking, int64, error)
}

type bookingService struct {
	bookingRepo repository.BookingRepository
	studioRepo  repository.StudioRepository
	slotRepo    repository.SlotRepository
	paymentRepo repository.PaymentRepository
	db          *gorm.DB
}

func NewBookingService(
	bookingRepo repository.BookingRepository,
	studioRepo repository.StudioRepository,
	slotRepo repository.SlotRepository,
	paymentRepo repository.PaymentRepository,
	db *gorm.DB,
) BookingService {
	return &bookingService{
		bookingRepo: bookingRepo,
		studioRepo:  studioRepo,
		slotRepo:    slotRepo,
		paymentRepo: paymentRepo,
		db:          db,
	}
}

func (s *bookingService) CreateBooking(ctx context.Context, input CreateBookingInput) (*domain.Booking, error) {
	if input.DurationHours < 1.0 {
		return nil, fmt.Errorf("%w: durasi minimal 1 jam", domain.ErrConflict)
	}

	// Hitung end time
	startParsed, err := time.Parse("15:04", input.StartTime)
	if err != nil {
		return nil, fmt.Errorf("format waktu tidak valid: gunakan HH:MM")
	}
	endParsed := startParsed.Add(time.Duration(float64(time.Hour) * input.DurationHours))
	endTime := endParsed.Format("15:04")

	// Validasi tanggal tidak di masa lalu
	today := time.Now().Truncate(24 * time.Hour)
	if input.BookingDate.Before(today) {
		return nil, fmt.Errorf("%w: tidak bisa booking di tanggal yang sudah lewat", domain.ErrConflict)
	}

	// Ambil data studio
	studio, err := s.studioRepo.FindByID(ctx, input.StudioID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	// Cek jam operasional
	dayName := weekdayToString(input.BookingDate.Weekday())
	slot, err := s.slotRepo.FindByStudioAndDay(ctx, input.StudioID, dayName)
	if err != nil || !slot.IsOpen {
		return nil, domain.ErrStudioClosed
	}

	openParsed, _  := time.Parse("15:04", slot.OpenTime[:5])
	closeParsed, _ := time.Parse("15:04", slot.CloseTime[:5])

	if startParsed.Before(openParsed) || endParsed.After(closeParsed) {
		return nil, domain.ErrOutsideOperatingHours
	}

	// Hitung biaya
	addonsCost := 0.0
	for _, a := range input.SelectedAddons {
		addonsCost += a.Price
	}
	serviceFee := studio.PricePerHour * input.DurationHours * 0.05 // 5%
	totalAmount := (studio.PricePerHour * input.DurationHours) + addonsCost + serviceFee

	// Serialisasi selected addons ke JSON
	addonsJSON, err := json.Marshal(input.SelectedAddons)
	if err != nil {
		return nil, err
	}

	var createdBooking *domain.Booking

	// BEGIN TRANSACTION — anti double-booking
	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. SELECT FOR UPDATE NOWAIT: lock baris yang overlap
		overlapping, err := s.bookingRepo.FindOverlapping(
			ctx, tx, input.StudioID, input.BookingDate,
			input.StartTime+":00", endTime+":00",
		)
		if err != nil {
			// Error lock (locked by another transaction) → slot tidak tersedia
			return domain.ErrSlotNotAvailable
		}
		if len(overlapping) > 0 {
			return domain.ErrSlotNotAvailable
		}

		// 2. Buat booking
		newBooking := &domain.Booking{
			ID:             uuid.New(),
			UserID:         input.UserID,
			StudioID:       input.StudioID,
			BookingDate:    input.BookingDate,
			StartTime:      input.StartTime + ":00",
			EndTime:        endTime + ":00",
			DurationHours:  input.DurationHours,
			Status:         domain.BookingStatusPending,
			PricePerHour:   studio.PricePerHour,
			AddonsCost:     addonsCost,
			ServiceFee:     serviceFee,
			TotalAmount:    totalAmount,
			SelectedAddons: addonsJSON,
			Notes:          input.Notes,
		}

		if err := s.bookingRepo.Create(ctx, tx, newBooking); err != nil {
			return err
		}

		// 3. Buat payment record (mock, expired 24 jam)
		expiredAt := time.Now().Add(24 * time.Hour)
		payment := &domain.Payment{
			ID:            uuid.New(),
			BookingID:     newBooking.ID,
			Amount:        totalAmount,
			PaymentMethod: domain.PaymentMethodMock,
			Status:        domain.PaymentStatusPending,
			ExpiredAt:     &expiredAt,
		}

		if err := s.paymentRepo.Create(ctx, tx, payment); err != nil {
			return err
		}

		newBooking.Payment = payment
		createdBooking = newBooking
		return nil // COMMIT
	})

	if txErr != nil {
		return nil, txErr
	}

	return createdBooking, nil
}

func (s *bookingService) GetUserBookings(ctx context.Context, userID uuid.UUID, status string, page, limit int) ([]domain.Booking, int64, error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 50 { limit = 10 }
	return s.bookingRepo.FindByUserID(ctx, userID, status, page, limit)
}

func (s *bookingService) GetBookingByID(ctx context.Context, bookingID, requesterID uuid.UUID, requesterRole string) (*domain.Booking, error) {
	booking, err := s.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	// Hanya pemilik booking atau studio_admin yang bisa akses
	if requesterRole != string(domain.RoleStudioAdmin) && booking.UserID != requesterID {
		return nil, domain.ErrForbidden
	}
	return booking, nil
}

func (s *bookingService) CancelBooking(ctx context.Context, bookingID, userID uuid.UUID) error {
	booking, err := s.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		return domain.ErrNotFound
	}
	if booking.UserID != userID {
		return domain.ErrForbidden
	}
	if booking.Status != domain.BookingStatusPending {
		return fmt.Errorf("%w: hanya booking dengan status 'pending' yang bisa dibatalkan", domain.ErrConflict)
	}
	return s.bookingRepo.UpdateStatus(ctx, bookingID, domain.BookingStatusCancelled)
}

func (s *bookingService) GetStudioBookings(ctx context.Context, studioID, adminID uuid.UUID, page, limit int) ([]domain.Booking, int64, error) {
	studio, err := s.studioRepo.FindByID(ctx, studioID)
	if err != nil {
		return nil, 0, domain.ErrNotFound
	}
	// Verifikasi admin adalah pengelola studio ini
	if studio.ManagedBy == nil || *studio.ManagedBy != adminID {
		return nil, 0, domain.ErrForbidden
	}
	if page < 1 { page = 1 }
	if limit < 1 || limit > 50 { limit = 10 }
	return s.bookingRepo.FindByStudioID(ctx, studioID, page, limit)
}
