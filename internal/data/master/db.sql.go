package master

import (
	"database/sql"

	"github.com/conneroisu/conneroh.com/internal/data/generic"
)

// New creates a new Queries instance.
func New(db generic.DBTX) *Queries {
	return &Queries{db: db}
}

// Queries is a wrapper around sql.DB that adds some convenience methods.
type Queries struct {
	db generic.DBTX
}

// WithTx returns a new Queries instance with the given transaction.
func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}
