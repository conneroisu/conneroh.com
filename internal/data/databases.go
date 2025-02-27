// Package data provides the data structures and methods for the data
// that is stored in the database. This package also provides the
// connection functionality to the database.
package data

import (
	"database/sql"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/conneroisu/conneroh.com/internal/data/generic"
	"github.com/conneroisu/conneroh.com/internal/data/master"

	// Register the libsql driver
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	// Register the modernc sqlite driver
	_ "modernc.org/sqlite"
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

// Config is a struct that holds the configuration for a database.
type Config struct {
	Schema   string
	URI      string
	FileName string
	Seed     string
}

// NewDb sets up the database using the URI and optional options.
// Using generics to return the appropriate type of query struct,
// it creates a new database struct with the sql database and the
// queries struct utilizing the URI and optional options provided.
func NewDb[
	Q master.Queries,
](
	newFunc func(generic.DBTX) *Q,
	config *Config,
) (*Database[Q], error) {
	u, err := url.Parse(config.URI)
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %v", err)
	}
	// if the file name contains a directory, ensure it exists (create if not)
	if strings.Contains(config.FileName, "/") {
		dir := filepath.Dir(config.FileName)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create directory: %v", err)
			}
		}
	}
	switch u.Scheme {
	case "file":
		db, err := sql.Open("sqlite", config.FileName)
		if err != nil {
			return nil,
				fmt.Errorf("failed to open db: %v", err)
		}
		if len(config.Schema) > 0 {
			for schem := range strings.SplitSeq(config.Schema, ";") {
				if len(strings.TrimSpace(schem)) == 0 {
					continue
				}
				_, err := db.Exec(schem)
				if err != nil {
					return nil,
						fmt.Errorf("error executing schema: %v, '%s'", err, schem)
				}
			}
		}
		if len(config.Seed) > 0 {
			for seed := range strings.SplitSeq(config.Seed, ";") {
				if len(strings.TrimSpace(seed)) == 0 {
					continue
				}
				_, err := db.Exec(seed)
				if err != nil {
					return nil,
						fmt.Errorf("error seeding db: \nseed: %s\nerror: %v", seed, err)
				}
			}
		}
		return &Database[Q]{
			Queries: newFunc(db),
			DB:      db,
		}, nil
	case "libsql":
		db, err := sql.Open("libsql", u.String())
		if err != nil {
			return nil,
				fmt.Errorf("failed to open db: %v", err)
		}
		if len(config.Schema) > 0 {
			for schem := range strings.SplitSeq(config.Schema, ";") {
				if len(strings.TrimSpace(schem)) == 0 {
					continue
				}
				_, err := db.Exec(schem)
				if err != nil {
					return nil,
						fmt.Errorf("error executing schema: %v, '%s'", err, schem)
				}
			}
		}
		if len(config.Seed) > 0 {
			for seed := range strings.SplitSeq(config.Seed, ";") {
				if len(strings.TrimSpace(seed)) == 0 {
					continue
				}
				if !strings.Contains(seed, "INSERT") {
					continue
				}
				_, err := db.Exec(seed)
				if err != nil {
					return nil,
						fmt.Errorf("error seeding db: \nseed: %s\nerror: %v", seed, err)
				}
			}
		}
		return &Database[Q]{
			Queries: newFunc(db),
			DB:      db,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}
}

// Close closes the database connection.
func (d *Database[Q]) Close() error {
	return d.DB.Close()
}
