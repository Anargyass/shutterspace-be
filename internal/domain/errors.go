package domain

import "errors"

// Domain-level errors — gunakan errors.Is() untuk perbandingan
var (
	ErrNotFound              = errors.New("resource not found")
	ErrConflict              = errors.New("resource conflict")
	ErrSlotNotAvailable      = errors.New("slot not available")
	ErrOutsideOperatingHours = errors.New("outside operating hours")
	ErrStudioClosed          = errors.New("studio closed on this day")
	ErrForbidden             = errors.New("access forbidden")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrEmailExists           = errors.New("email already exists")
	ErrInactive              = errors.New("account is inactive")
)
