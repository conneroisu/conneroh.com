package master_test

import (
	"testing"

	"github.com/conneroisu/conneroh.com/internal/data"
	"github.com/conneroisu/conneroh.com/internal/data/master"

	_ "embed"
)

//go:embed combined/schema.sql
var schema string

func TestSchema(t *testing.T) {
	db, err := data.NewDb(master.New, &data.Config{
		Schema:   schema,
		FileName: ":memory:",
		URI:      "file://dummy",
	})
	if err != nil {
		return
	}
	defer db.Close()
}
