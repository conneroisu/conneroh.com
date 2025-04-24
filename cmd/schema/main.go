package main

import (
	"context"
	"database/sql"

	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"

	_ "modernc.org/sqlite"
)

func main() {
	if err := run(context.Background()); err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {
	sqldb, err := sql.Open("sqlite", assets.DBName())
	if err != nil {
		return err
	}
	defer sqldb.Close()
	db := bun.NewDB(sqldb, sqlitedialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	err = assets.InitDB(ctx, db)
	if err != nil {
		return err
	}

	return nil
}
