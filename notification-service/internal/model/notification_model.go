package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Notifications struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null" json:"order_id"`
	Type      string    `gorm:"type:varchar(50);not null" json:"type"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	IsRead    bool      `gorm:"type:bool;default:false" json:"is_read"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}