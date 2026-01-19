package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Orders struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Status      string `gorm:"type:varchar(50);default:'pending';not null" json:"status"`
	Email       string `gorm:"type:varchar(255);not null" json:"email"`
	TotalAmount float64 `gorm:"type:numeric(12,2);not null" json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
	PaymentMethod string `gorm:"type:varchar(50);" json:"payment_method"`
	PaymentID 		*uuid.UUID `gorm:"type:uuid;" json:"payment_id"`
	PaymentStatus string `gorm:"type:varchar(50);default:'PENDING';" json:"payment_status"`
}