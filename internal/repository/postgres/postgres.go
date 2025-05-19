package postgres

import (
	"fmt"
	"log"

	"github.com/imranzahoor/banking-ledger/internal/model"
	"github.com/imranzahoor/banking-ledger/pkg/config"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(cfg config.Config) (*gorm.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPass,
		cfg.PostgresDB,
	)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.Account{}, &model.Transaction{})
	if err != nil {
		return nil, err
	}

	log.Println("Postgres connected and migrations applied.")
	return db, nil
}
