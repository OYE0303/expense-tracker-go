package model

import (
	"context"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionModel struct {
	DB *mongo.Database
}

// TODO add a type for the model

func newTransactionModel(db *mongo.Database) *TransactionModel {
	return &TransactionModel{DB: db}
}

func (t *TransactionModel) Create(ctx context.Context, transaction *domain.Transaction) error {
	// TODO set created_at and updated_at manually
	// TODO convert domain.Transaction to model.Transaction(type)

	_, err := t.DB.Collection("transactions").InsertOne(ctx, transaction)
	if err != nil {
		return err
	}
	return nil
}
