package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string
type PaymentMethod string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"

	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodCard         PaymentMethod = "credit_debit_card"
	PaymentMethodEWallet      PaymentMethod = "e_wallet"
	PaymentMethodMock         PaymentMethod = "mock_payment"
)

type Payment struct {
	ID            uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BookingID     uuid.UUID     `gorm:"type:uuid;not null;uniqueIndex"                 json:"booking_id"`
	Amount        float64       `gorm:"type:numeric(12,2);not null"                    json:"amount"`
	PaymentMethod PaymentMethod `gorm:"type:payment_method;default:'mock_payment'"     json:"payment_method"`
	Status        PaymentStatus `gorm:"type:payment_status;default:'pending'"          json:"status"`
	ExternalRef   string        `gorm:"type:varchar(100)"                              json:"external_ref"`
	PaidAt        *time.Time    `                                                      json:"paid_at"`
	ExpiredAt     *time.Time    `                                                      json:"expired_at"`
	CreatedAt     time.Time     `                                                      json:"created_at"`
	UpdatedAt     time.Time     `                                                      json:"updated_at"`
}
