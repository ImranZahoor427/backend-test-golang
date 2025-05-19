package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/imranzahoor/banking-ledger/pkg/constants"
)

// Transaction represents a banking ledger entry.
type Transaction struct {
	ID          uuid.UUID                 `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	AccountID   uuid.UUID                 `gorm:"type:uuid;not null"`
	Type        constants.TransactionType `gorm:"type:varchar(20);not null"`
	Amount      int64                     `gorm:"not null"`
	Description string
	CreatedAt   time.Time
}
