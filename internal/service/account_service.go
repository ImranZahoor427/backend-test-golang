package service

import (
	"context"
	"errors"
	"time"

	"github.com/imranzahoor/banking-ledger/internal/model"
	"github.com/imranzahoor/banking-ledger/internal/repository/postgres"
)

type AccountServiceInterface interface {
	CreateAccount(ctx context.Context, ownerName string, initialBalance int64) (*model.Account, error)
	GetAccountByID(ctx context.Context, accountID string) (*model.Account, error)
}
type AccountService struct {
	accountRepo postgres.AccountRepository
}

func NewAccountService(repo postgres.AccountRepository) *AccountService {
	return &AccountService{accountRepo: repo}
}

func (s *AccountService) CreateAccount(ctx context.Context, ownerName string, initialBalance int64) (*model.Account, error) {
	if ownerName == "" {
		return nil, errors.New("owner name required")
	}
	if initialBalance < 0 {
		return nil, errors.New("initial balance cannot be negative")
	}

	acc := &model.Account{
		OwnerName: ownerName,
		Balance:   initialBalance,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err := s.accountRepo.CreateAccount(ctx, acc)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (s *AccountService) GetAccountByID(ctx context.Context, accountID string) (*model.Account, error) {
	return s.accountRepo.GetAccountByID(ctx, accountID)
}
