package entity

import (
	"time"

	uuid "github.com/jackc/pgx/pgtype/ext/gofrs-uuid"
)

type Notifications struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrderId uuid.UUID 
	Type string 
	Message string 
	IsRead bool 
	CreatedAt time.Time
}