package entity

import (
	"time"

	uuid "github.com/jackc/pgx/pgtype/ext/gofrs-uuid"
)

type Orders struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Status string 
	Email string
	TotalAmount float64 
	CreatedAt time.Time
}