package queue

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/imranzahoor/banking-ledger/internal/model"
	"github.com/imranzahoor/banking-ledger/internal/repository/mongo"
	"github.com/imranzahoor/banking-ledger/internal/repository/postgres"
	"github.com/imranzahoor/banking-ledger/pkg/config"
	"github.com/imranzahoor/banking-ledger/pkg/constants"
	"github.com/segmentio/kafka-go"
)

type TransactionConsumer struct {
	kafkaReader *kafka.Reader
	accountRepo *postgres.AccountRepo
	ledgerRepo  *mongo.LedgerRepo
}

func NewTransactionConsumer(cfg config.KafkaConfig, ar *postgres.AccountRepo, lr *mongo.LedgerRepo) *TransactionConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		GroupID:  cfg.GroupID,
		Topic:    cfg.Topic,
		MinBytes: 1e3,  // 1KB
		MaxBytes: 10e6, // 10MB
	})
	return &TransactionConsumer{
		kafkaReader: r,
		accountRepo: ar,
		ledgerRepo:  lr,
	}
}

func (c *TransactionConsumer) Run(ctx context.Context) error {
	for {
		m, err := c.kafkaReader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		var txn model.Transaction
		if err := json.Unmarshal(m.Value, &txn); err != nil {
			log.Printf("invalid transaction payload: %v", err)
			continue
		}

		// Process transaction (reuse ProcessTransaction logic)
		err = c.processTransaction(ctx, &txn)
		if err != nil {
			log.Printf("failed to process transaction ID %s: %v", txn.ID, err)
		}
	}
}

func (c *TransactionConsumer) processTransaction(ctx context.Context, txn *model.Transaction) error {
	var delta int64
	switch txn.Type {
	case constants.Deposit:
		delta = txn.Amount
	case constants.Withdrawal:
		delta = -txn.Amount
	default:
		return errors.New("invalid transaction type")
	}

	if err := c.accountRepo.UpdateBalance(ctx, txn.AccountID, delta); err != nil {
		return err
	}

	if err := c.ledgerRepo.InsertTransaction(ctx, txn); err != nil {
		// rollback balance update
		_ = c.accountRepo.UpdateBalance(ctx, txn.AccountID, -delta)
		return err
	}

	return nil
}
