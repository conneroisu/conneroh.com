// Package main is the main package for the updated the generated code.
package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/conneroisu/conneroh.com/internal/cache"
	"github.com/conneroisu/conneroh.com/internal/llama"
	"github.com/conneroisu/conneroh.com/internal/logger"
	"github.com/conneroisu/conneroh.com/internal/tigris"
	"github.com/spf13/afero"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"
	_ "modernc.org/sqlite"
)

var (
	workers = flag.Int("jobs", 80, "number of parallel workers")
)

func main() {
	flag.Parse()
	slog.SetDefault(logger.DefaultLogger)
	if err := run(
		context.Background(),
		os.Getenv,
		*workers,
	); err != nil {
		panic(err)
	}
}

func run(
	ctx context.Context,
	getenv func(string) string,
	workers int,
) error {
	var (
		errCh = make(chan error)
		msgCh = make(chan *cache.Msg)
		wg    = sync.WaitGroup{}
	)
	innerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	ol, err := llama.NewOllamaClient(getenv)
	if err != nil {
		return err
	}
	ti, err := tigris.New(getenv)
	if err != nil {
		return err
	}
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
	fs := afero.NewBasePathFs(afero.NewOsFs(), assets.VaultLoc)
	md := assets.NewMD(fs)

	h := cache.NewHollywood(workers, fs, db, ti, ol, md)
	println("starting hollywood")

	go h.Start(innerCtx, msgCh, errCh, &wg)
	go ReadFS(fs, &wg, msgCh, errCh)
	time.Sleep(time.Millisecond * 50)
	go func() {
		wg.Wait()
		cancel()
	}()
	for {
		select {
		case <-time.After(time.Second * 1):
			slog.Info("waiting for hollywood to finish")
		case <-innerCtx.Done():
			close(msgCh)
			close(errCh)
			return nil
		case err := <-errCh:
			cancel()
			close(msgCh)
			close(errCh)
			return err
		}
	}
}

// ReadFS reads the filesystem and sends messages to the message channel.
func ReadFS(
	fs afero.Fs,
	wg *sync.WaitGroup,
	msgCh chan *cache.Msg,
	errCh chan error,
) {
	wg.Add(1)
	defer wg.Done()
	err := afero.Walk(
		fs,
		".",
		func(fPath string, info os.FileInfo, err error) error {
			wg.Add(1)
			defer wg.Done()

			var msgType cache.MsgType
			switch {
			case err != nil:
				return err
			case assets.IsAllowedMediaType(fPath):
				msgType = cache.MsgTypeAsset
			case assets.IsAllowedDocumentType(fPath):
				msgType = cache.MsgTypeDoc
			case info.IsDir():
				return nil
			default:
				return nil
			}
			msgCh <- &cache.Msg{Path: fPath, Type: msgType}
			return nil
		})
	if err != nil {
		errCh <- err
	}
}
