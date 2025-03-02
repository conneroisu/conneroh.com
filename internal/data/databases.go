// Package data provides the data structures and methods for the data
// that is stored in the database. This package also provides the
// connection functionality to the database.
package data

import (
	"database/sql"
	"fmt"
	"io"
	"net/url"

	"github.com/conneroisu/conneroh.com/internal/data/generic"
	"github.com/conneroisu/conneroh.com/internal/data/master"

	// Register the libsql driver
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var _ io.Closer = (*Database[master.Queries])(nil)

// Database is a struct that holds the sql database and the queries.
// It uses generics to hold the appropriate type of query struct.
type Database[
	T master.Queries,
] struct {
	Queries *T
	DB      *sql.DB
}

// NewDb sets up the database using the URI and optional options.
//
// Using generics to return the appropriate type of query struct,
// it creates a new database struct with the sql database and the
// queries struct utilizing the URI and optional options provided.
func NewDb[
	Q master.Queries,
](
	newFunc func(generic.DBTX) *Q,
	URI string,
) (*Database[Q], error) {
	u, err := url.Parse(URI)
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %v", err)
	}
	db, err := sql.Open("libsql", u.String())
	if err != nil {
		return nil,
			fmt.Errorf("failed to open db: %v", err)
	}
	return &Database[Q]{
		Queries: newFunc(db),
		DB:      db,
	}, nil
}

// Close closes the database connection.
func (d *Database[Q]) Close() error {
	return d.DB.Close()
}
