package model

import (
	"time"

	"github.com/google/uuid"
)

// Account represents a bank account domain model.
type Account struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OwnerName string    `gorm:"not null"`
	Balance   int64     `gorm:"not null"` // smallest currency unit (e.g. cents)
	CreatedAt time.Time
	UpdatedAt time.Time
}
