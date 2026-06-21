package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"shutterspace/internal/domain"
	"shutterspace/internal/repository"
)

type AvailabilityResult struct {
	Date           string       `json:"date"`
	Day            string       `json:"day"`
	StudioID       string       `json:"studio_id"`
	OperatingHours OpHours      `json:"operating_hours"`
	BookedSlots    []BookedSlot `json:"booked_slots"`
}

type OpHours struct {
	Open  string `json:"open"`
	Close string `json:"close"`
}

type BookedSlot struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type AvailabilityService interface {
	GetAvailability(ctx context.Context, studioID uuid.UUID, date time.Time) (*AvailabilityResult, error)
}

type availabilityService struct {
	slotRepo    repository.SlotRepository
	bookingRepo repository.BookingRepository
}

func NewAvailabilityService(slotRepo repository.SlotRepository, bookingRepo repository.BookingRepository) AvailabilityService {
	return &availabilityService{
		slotRepo:    slotRepo,
		bookingRepo: bookingRepo,
	}
}

func (s *availabilityService) GetAvailability(ctx context.Context, studioID uuid.UUID, date time.Time) (*AvailabilityResult, error) {
	dayName := weekdayToString(date.Weekday())

	slot, err := s.slotRepo.FindByStudioAndDay(ctx, studioID, dayName)
	if err != nil {
		return nil, domain.ErrStudioClosed
	}

	if !slot.IsOpen {
		return nil, domain.ErrStudioClosed
	}

	bookings, err := s.bookingRepo.FindByStudioAndDate(ctx, studioID, date)
	if err != nil {
		return nil, err
	}

	bookedSlots := make([]BookedSlot, 0, len(bookings))
	for _, b := range bookings {
		start := b.StartTime
		end := b.EndTime
		// Normalize time format ke "HH:MM"
		if len(start) >= 5 {
			start = start[:5]
		}
		if len(end) >= 5 {
			end = end[:5]
		}
		bookedSlots = append(bookedSlots, BookedSlot{Start: start, End: end})
	}

	// Normalize operating hours ke "HH:MM"
	openTime := slot.OpenTime
	closeTime := slot.CloseTime
	if len(openTime) >= 5 {
		openTime = openTime[:5]
	}
	if len(closeTime) >= 5 {
		closeTime = closeTime[:5]
	}

	return &AvailabilityResult{
		Date:     date.Format("2006-01-02"),
		Day:      dayName,
		StudioID: studioID.String(),
		OperatingHours: OpHours{
			Open:  openTime,
			Close: closeTime,
		},
		BookedSlots: bookedSlots,
	}, nil
}

func weekdayToString(d time.Weekday) string {
	m := map[time.Weekday]string{
		time.Monday:    "monday",
		time.Tuesday:   "tuesday",
		time.Wednesday: "wednesday",
		time.Thursday:  "thursday",
		time.Friday:    "friday",
		time.Saturday:  "saturday",
		time.Sunday:    "sunday",
	}
	return m[d]
}
