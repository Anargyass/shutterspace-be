package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"shutterspace/internal/domain"
	"shutterspace/internal/service"
	"shutterspace/pkg/response"
)

type AvailabilityHandler struct {
	availabilityService service.AvailabilityService
}

func NewAvailabilityHandler(as service.AvailabilityService) *AvailabilityHandler {
	return &AvailabilityHandler{availabilityService: as}
}

func (h *AvailabilityHandler) GetAvailability(c *gin.Context) {
	studioIDStr := c.Param("id")
	studioID, err := uuid.Parse(studioIDStr)
	if err != nil {
		response.BadRequest(c, "INVALID_ID", "Studio ID tidak valid")
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		response.BadRequest(c, "MISSING_DATE", "Parameter 'date' (YYYY-MM-DD) wajib diisi")
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.BadRequest(c, "INVALID_DATE", "Format tanggal harus YYYY-MM-DD")
		return
	}

	avail, err := h.availabilityService.GetAvailability(c.Request.Context(), studioID, date)
	if err != nil {
		switch err {
		case domain.ErrStudioClosed:
			response.NotFound(c, "Studio tutup pada hari tersebut")
		default:
			response.InternalError(c, "Gagal mengambil data ketersediaan")
		}
		return
	}

	response.OK(c, avail)
}
