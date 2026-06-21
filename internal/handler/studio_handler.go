package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"shutterspace/internal/repository"
	"shutterspace/internal/service"
	"shutterspace/pkg/response"
)

type StudioHandler struct {
	studioService service.StudioService
}

func NewStudioHandler(ss service.StudioService) *StudioHandler {
	return &StudioHandler{studioService: ss}
}

func (h *StudioHandler) ListStudios(c *gin.Context) {
	page, _  := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)

	filter := repository.StudioFilter{
		TypeSlug: c.Query("type"),
		Area:     c.Query("area"),
		MinPrice: minPrice,
		MaxPrice: maxPrice,
		Page:     page,
		Limit:    limit,
	}

	studios, total, err := h.studioService.ListStudios(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, "Gagal mengambil data studio")
		return
	}

	response.OK(c, gin.H{
		"studios": studios,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func (h *StudioHandler) GetStudio(c *gin.Context) {
	idStr := c.Param("id")

	// Coba parse sebagai UUID dulu, kalau gagal coba sebagai slug
	id, err := uuid.Parse(idStr)
	if err != nil {
		// Anggap sebagai slug
		studio, err := h.studioService.GetStudioBySlug(c.Request.Context(), idStr)
		if err != nil {
			response.NotFound(c, "Studio tidak ditemukan")
			return
		}
		response.OK(c, gin.H{"studio": studio})
		return
	}

	studio, err := h.studioService.GetStudioByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "Studio tidak ditemukan")
		return
	}

	response.OK(c, gin.H{"studio": studio})
}

func (h *StudioHandler) GetStudioTypes(c *gin.Context) {
	types, err := h.studioService.GetAllStudioTypes(c.Request.Context())
	if err != nil {
		response.InternalError(c, "Gagal mengambil tipe studio")
		return
	}
	response.OK(c, gin.H{"types": types})
}

func (h *StudioHandler) GetMyStudios(c *gin.Context) {
	adminIDStr := c.GetString("user_id")
	adminID, _ := uuid.Parse(adminIDStr)

	studios, err := h.studioService.GetMyStudios(c.Request.Context(), adminID)
	if err != nil {
		response.InternalError(c, "Gagal mengambil data studio")
		return
	}

	response.OK(c, gin.H{"studios": studios})
}
