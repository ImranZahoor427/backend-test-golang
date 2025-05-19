package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/imranzahoor/banking-ledger/internal/api"
	"github.com/imranzahoor/banking-ledger/internal/middleware"
	"github.com/imranzahoor/banking-ledger/internal/repository/mongo"
	"github.com/imranzahoor/banking-ledger/internal/repository/postgres"
	"github.com/imranzahoor/banking-ledger/internal/service"
	"github.com/imranzahoor/banking-ledger/pkg/config"
	queue "github.com/imranzahoor/banking-ledger/pkg/kafka"
)

func main() {
	cfg := config.Load()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	accountRepo := initPostgresAccountRepo(cfg)
	transactionRepo := initMongoLedgerRepo(cfg)

	startTransactionConsumer(ctx, cfg, accountRepo, transactionRepo)

	accountService := service.NewAccountService(accountRepo)
	transactionService := service.NewTransactionService(accountRepo, transactionRepo, cfg)

	runHTTPServer(cfg, accountService, transactionService)
}

func initPostgresAccountRepo(cfg config.Config) *postgres.AccountRepo {
	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to Postgres: %v", err)
	}
	return postgres.NewAccountRepo(db)
}

func initMongoLedgerRepo(cfg config.Config) *mongo.LedgerRepo {
	client, err := mongo.NewMongoClient(cfg)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	return mongo.NewLedgerRepo(client, cfg.MongoDB)
}

func startTransactionConsumer(ctx context.Context, cfg config.Config, ar *postgres.AccountRepo, lr *mongo.LedgerRepo) {
	kafkaCfg := config.KafkaConfig{
		Brokers: cfg.KafkaBrokers,
		Topic:   cfg.KafkaTopic,
		GroupID: cfg.KafkaGroupID,
	}
	consumer := queue.NewTransactionConsumer(kafkaCfg, ar, lr)
	go func() {
		if err := consumer.Run(ctx); err != nil {
			log.Fatalf("Kafka consumer failed: %v", err)
		}
	}()
}

func runHTTPServer(cfg config.Config, accountService *service.AccountService, transactionService *service.TransactionService) {
	router := gin.Default()
	router.Use(middleware.Recovery())

	handler := api.NewHandler(accountService, transactionService)
	handler.RegisterRoutes(router)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
