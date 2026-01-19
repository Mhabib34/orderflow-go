package model

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusSuccess PaymentStatus = "success"
	PaymentStatusFailed  PaymentStatus = "failed"
	PaymentStatusExpired PaymentStatus = "expired"
)

type Payments struct {
	ID      uuid.UUID 		`gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrderID uuid.UUID 		`gorm:"type:uuid;not null" json:"order_id"`
	Amount  int64     		`gorm:"type:bigint;not null" json:"amount"`
	Method  string    		`gorm:"type:varchar(50);not null" json:"method"`
	ProviderRefID string 	`gorm:"type:varchar(255);not null" json:"provider_ref_id"`
	Status  PaymentStatus   `gorm:"type:varchar(50);default:'pending';not null" json:"status"`
	PaidAt  time.Time 		`json:"paid_at"`
	ExpiredAt time.Time 	`json:"expired_at"`
	CreatedAt time.Time 	`gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time 	`gorm:"autoUpdateTime" json:"updated_at"`
	PaymentURL string 		`gorm:"type:varchar(255);not null" json:"payment_url"`
	PaymentID  string 		`gorm:"type:varchar(255);not null" json:"payment_id"`
}