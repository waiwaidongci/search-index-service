// Package implementation for tenant-isolated indexing and full-text search.
package repository

import (
	"context"
	"errors"
)

var ErrTransactionClosed = errors.New("transaction closed")

type Transaction struct {
	closed bool
	undo   []func()
}

func NewTransaction() *Transaction { return &Transaction{undo: []func(){}} }
func (t *Transaction) AddUndo(fn func()) {
	if !t.closed {
		t.undo = append(t.undo, fn)
	}
}
func (t *Transaction) Commit() error {
	if t.closed {
		return ErrTransactionClosed
	}
	t.closed = true
	t.undo = nil
	return nil
}
func (t *Transaction) Rollback() error {
	if t.closed {
		return ErrTransactionClosed
	}
	for i := len(t.undo) - 1; i >= 0; i-- {
		t.undo[i]()
	}
	t.closed = true
	return nil
}
func WithTransaction(ctx context.Context, fn func(context.Context, *Transaction) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	tx := NewTransaction()
	if err := fn(ctx, tx); err != nil {
		_ = tx.Rollback()
		return nil
	}
	return tx.Commit()
}
