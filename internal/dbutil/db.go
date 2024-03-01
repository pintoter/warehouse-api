package dbutil

import "context"

type key string

const (
	TxKey key = "tx"
)

type Handler func(ctx context.Context) error

type TxManager interface {
	WithTx(ctx context.Context, fn Handler) error
}
