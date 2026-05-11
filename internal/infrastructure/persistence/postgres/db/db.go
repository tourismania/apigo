// Package db hosts the sqlc-shaped query layer. The file layout mirrors
// what `sqlc generate` would emit (db.go / models.go / users.sql.go), so
// regenerating against sqlc.yaml is a drop-in replacement. The code is
// hand-written today to keep the project buildable without an extra
// codegen step in CI.
package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DBTX is the minimum surface of pgx the generated queries need. It is
// implemented by both *pgxpool.Pool and pgx.Tx, so callers can wrap
// queries in a transaction transparently.
type DBTX interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

// Queries is the sqlc-style aggregate of all generated query methods.
type Queries struct {
	db DBTX
}

// New constructs a Queries over the given pool/tx.
func New(db DBTX) *Queries {
	return &Queries{db: db}
}

// WithTx returns a copy of Queries bound to the given transaction. The
// caller owns the tx lifecycle.
func (q *Queries) WithTx(tx pgx.Tx) *Queries {
	return &Queries{db: tx}
}
