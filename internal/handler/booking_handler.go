package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"shutterspace/internal/domain"
	"shutterspace/internal/service"
	"shutterspace/pkg/response"
)

type BookingHandler struct {
	bookingService service.BookingService
}

func NewBookingHandler(bs service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bs}
}

type createBookingRequest struct {
	StudioID       string               `json:"studio_id"       binding:"required,uuid"`
	BookingDate    string               `json:"booking_date"    binding:"required"`
	StartTime      string               `json:"start_time"      binding:"required"`
	DurationHours  float64              `json:"duration_hours"  binding:"required,min=1"`
	SelectedAddons []domain.AddonItem   `json:"selected_addons"`
	Notes          string               `json:"notes"           binding:"max=500"`
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req createBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}

	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(c, "Token tidak valid")
		return
	}

	bookingDate, err := time.Parse("2006-01-02", req.BookingDate)
	if err != nil {
		response.BadRequest(c, "INVALID_DATE", "Format tanggal harus YYYY-MM-DD")
		return
	}

	studioID, _ := uuid.Parse(req.StudioID)

	booking, err := h.bookingService.CreateBooking(c.Request.Context(), service.CreateBookingInput{
		UserID:         userID,
		StudioID:       studioID,
		BookingDate:    bookingDate,
		StartTime:      req.StartTime,
		DurationHours:  req.DurationHours,
		SelectedAddons: req.SelectedAddons,
		Notes:          req.Notes,
	})

	if err != nil {
		switch err {
		case domain.ErrSlotNotAvailable:
			response.Conflict(c, "SLOT_NOT_AVAILABLE", "Waktu tersebut sudah dipesan atau sedang dalam proses pemesanan")
		case domain.ErrOutsideOperatingHours:
			response.BadRequest(c, "OUTSIDE_OPERATING_HOURS", "Booking melebihi jam operasional studio")
		case domain.ErrStudioClosed:
			response.BadRequest(c, "STUDIO_CLOSED", "Studio tutup pada hari tersebut")
		case domain.ErrNotFound:
			response.NotFound(c, "Studio tidak ditemukan")
		default:
			response.InternalError(c, "Gagal membuat booking")
		}
		return
	}

	response.Created(c, gin.H{"booking": booking})
}

func (h *BookingHandler) GetMyBookings(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)

	status := c.Query("status")
	page, _  := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	bookings, total, err := h.bookingService.GetUserBookings(c.Request.Context(), userID, status, page, limit)
	if err != nil {
		response.InternalError(c, "Gagal mengambil data booking")
		return
	}

	response.OK(c, gin.H{
		"bookings": bookings,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (int(total) + limit - 1) / limit,
		},
	})
}

func (h *BookingHandler) GetBookingByID(c *gin.Context) {
	bookingIDStr := c.Param("id")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		response.BadRequest(c, "INVALID_ID", "Booking ID tidak valid")
		return
	}

	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)
	userRole := c.GetString("user_role")

	booking, err := h.bookingService.GetBookingByID(c.Request.Context(), bookingID, userID, userRole)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			response.NotFound(c, "Booking tidak ditemukan")
		case domain.ErrForbidden:
			response.Forbidden(c, "Anda tidak memiliki akses ke booking ini")
		default:
			response.InternalError(c, "Gagal mengambil data booking")
		}
		return
	}

	response.OK(c, gin.H{"booking": booking})
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	bookingIDStr := c.Param("id")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		response.BadRequest(c, "INVALID_ID", "Booking ID tidak valid")
		return
	}

	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)

	if err := h.bookingService.CancelBooking(c.Request.Context(), bookingID, userID); err != nil {
		switch err {
		case domain.ErrNotFound:
			response.NotFound(c, "Booking tidak ditemukan")
		case domain.ErrForbidden:
			response.Forbidden(c, "Anda tidak memiliki akses ke booking ini")
		default:
			response.Conflict(c, "CANCEL_FAILED", err.Error())
		}
		return
	}

	response.OK(c, gin.H{"message": "Booking berhasil dibatalkan"})
}

func (h *BookingHandler) GetStudioBookings(c *gin.Context) {
	studioIDStr := c.Param("studioId")
	studioID, err := uuid.Parse(studioIDStr)
	if err != nil {
		response.BadRequest(c, "INVALID_ID", "Studio ID tidak valid")
		return
	}

	adminIDStr := c.GetString("user_id")
	adminID, _ := uuid.Parse(adminIDStr)

	page, _  := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	bookings, total, err := h.bookingService.GetStudioBookings(c.Request.Context(), studioID, adminID, page, limit)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			response.NotFound(c, "Studio tidak ditemukan")
		case domain.ErrForbidden:
			response.Forbidden(c, "Anda bukan pengelola studio ini")
		default:
			response.InternalError(c, "Gagal mengambil data booking")
		}
		return
	}

	response.OK(c, gin.H{
		"bookings": bookings,
		"pagination": gin.H{"page": page, "limit": limit, "total": total},
	})
}
