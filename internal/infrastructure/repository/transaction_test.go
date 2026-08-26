package repository

import (
	"context"
	"errors"
	"testing"
)

func TestWithTransactionDoesNotSwallowRollbackError(t *testing.T) {
	err := WithTransaction(context.Background(), func(ctx context.Context, tx *Transaction) error {
		tx.AddUndo(func() {})
		return errors.New("business failure")
	})
	if err == nil {
		t.Fatal("expected WithTransaction to return business failure")
	}
}
