package service

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/imranzahoor/banking-ledger/internal/model"
	"github.com/imranzahoor/banking-ledger/internal/repository/mongo"
	"github.com/imranzahoor/banking-ledger/internal/repository/postgres"
	"github.com/imranzahoor/banking-ledger/pkg/config"
	"github.com/segmentio/kafka-go"
)

type TransactionService struct {
	accountRepo postgres.AccountRepository
	ledgerRepo  mongo.LedgerRepository
	kafkaWriter *kafka.Writer
}

type TransactionServiceInterface interface {
	EnqueueTransaction(ctx context.Context, txn *model.Transaction) error
	GetTransactions(accountID uuid.UUID, limit, offset int64) ([]model.Transaction, error)
	HasSufficientFunds(ctx context.Context, accountID string, amount int64) (bool, error)
}

func NewTransactionService(
	ar postgres.AccountRepository,
	lr mongo.LedgerRepository,
	cfg config.Config,
) *TransactionService {
	w := &kafka.Writer{
		Addr:     kafka.TCP(cfg.KafkaBrokers...),
		Topic:    cfg.KafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}

	return &TransactionService{
		accountRepo: ar,
		ledgerRepo:  lr,
		kafkaWriter: w,
	}
}

func (s *TransactionService) EnqueueTransaction(ctx context.Context, txn *model.Transaction) error {
	if txn.ID == uuid.Nil {
		txn.ID = uuid.New()
	}

	data, err := json.Marshal(txn)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(txn.ID.String()),
		Value: data,
	}

	return s.kafkaWriter.WriteMessages(ctx, msg)
}

func (s *TransactionService) GetTransactions(accountID uuid.UUID, limit, offset int64) ([]model.Transaction, error) {
	return s.ledgerRepo.GetTransactionsByAccountID(context.Background(), accountID, limit, offset)
}

func (s *TransactionService) HasSufficientFunds(ctx context.Context, accountID string, amount int64) (bool, error) {
	account, err := s.accountRepo.GetAccountByID(ctx, accountID)
	if err != nil {
		return false, err
	}
	return account.Balance >= amount, nil
}
