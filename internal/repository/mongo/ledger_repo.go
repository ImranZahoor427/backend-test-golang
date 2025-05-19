package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/google/uuid"
	"github.com/imranzahoor/banking-ledger/internal/model"
)

type LedgerRepo struct {
	coll *mongo.Collection
}
type LedgerRepository interface {
	InsertTransaction(ctx context.Context, txn *model.Transaction) error
	GetTransactionsByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int64) ([]model.Transaction, error)
}

func NewLedgerRepo(client *mongo.Client, dbName string) *LedgerRepo {
	coll := client.Database(dbName).Collection("transactions")
	return &LedgerRepo{coll: coll}
}

// InsertTransaction adds a new transaction log entry
func (r *LedgerRepo) InsertTransaction(ctx context.Context, txn *model.Transaction) error {
	if txn.ID == uuid.Nil {
		txn.ID = uuid.New()
	}

	txn.CreatedAt = time.Now().UTC()

	_, err := r.coll.InsertOne(ctx, txn)
	return err
}

// GetTransactionsByAccountID fetches transaction logs for account with optional limit/offset
func (r *LedgerRepo) GetTransactionsByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int64) ([]model.Transaction, error) {
	filter := bson.M{"accountid": accountID}

	findOptions := options.Find()
	if limit > 0 {
		findOptions.SetLimit(limit)
	}
	if offset > 0 {
		findOptions.SetSkip(offset)
	}
	findOptions.SetSort(bson.D{{Key: "createdat", Value: -1}}) // newest first

	cursor, err := r.coll.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []model.Transaction
	for cursor.Next(ctx) {
		var txn model.Transaction
		if err := cursor.Decode(&txn); err != nil {
			return nil, err
		}
		results = append(results, txn)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
