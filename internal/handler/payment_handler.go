package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"shutterspace/internal/domain"
	"shutterspace/internal/service"
	"shutterspace/pkg/response"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(ps service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: ps}
}

type processPaymentRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required,oneof=bank_transfer credit_debit_card e_wallet mock_payment"`
}

func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	bookingIDStr := c.Param("bookingId")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		response.BadRequest(c, "INVALID_ID", "Booking ID tidak valid")
		return
	}

	var req processPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}

	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)

	var method domain.PaymentMethod
	switch req.PaymentMethod {
	case "bank_transfer":
		method = domain.PaymentMethodBankTransfer
	case "credit_debit_card":
		method = domain.PaymentMethodCard
	case "e_wallet":
		method = domain.PaymentMethodEWallet
	default:
		method = domain.PaymentMethodMock
	}

	payment, err := h.paymentService.ProcessMockPayment(c.Request.Context(), bookingID, userID, method)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			response.NotFound(c, "Booking tidak ditemukan")
		case domain.ErrForbidden:
			response.Forbidden(c, "Anda tidak memiliki akses ke booking ini")
		default:
			response.Conflict(c, "PAYMENT_FAILED", err.Error())
		}
		return
	}

	response.OK(c, gin.H{
		"payment": payment,
		"message": "Pembayaran berhasil diproses",
	})
}
