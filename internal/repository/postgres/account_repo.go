package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/imranzahoor/banking-ledger/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountRepo struct {
	db *gorm.DB
}

type AccountRepository interface {
	CreateAccount(ctx context.Context, acc *model.Account) error
	GetAccountByID(ctx context.Context, id string) (*model.Account, error)
	UpdateBalance(ctx context.Context, accountID uuid.UUID, delta int64) error
}

func NewAccountRepo(db *gorm.DB) *AccountRepo {
	return &AccountRepo{db: db}
}

// CreateAccount inserts new account with initial balance
func (r *AccountRepo) CreateAccount(ctx context.Context, acc *model.Account) error {
	if acc.ID == uuid.Nil {
		acc.ID = uuid.New()
	}

	acc.CreatedAt = time.Now().UTC()
	acc.UpdatedAt = acc.CreatedAt

	return r.db.WithContext(ctx).Create(acc).Error
}

// GetAccountByID fetches account by ID
func (r *AccountRepo) GetAccountByID(ctx context.Context, id string) (*model.Account, error) {
	var acc model.Account
	err := r.db.WithContext(ctx).First(&acc, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

// UpdateBalance atomically updates the balance by delta amount (positive or negative)
// Returns error if balance would go negative.
func (r *AccountRepo) UpdateBalance(ctx context.Context, accountID uuid.UUID, delta int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var acc model.Account

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&acc, "id = ?", accountID.String()).Error; err != nil {
			return err
		}

		newBalance := acc.Balance + delta
		if newBalance < 0 {
			return errors.New("insufficient funds")
		}

		acc.Balance = newBalance
		acc.UpdatedAt = time.Now().UTC()
		return tx.Save(&acc).Error
	})
}
