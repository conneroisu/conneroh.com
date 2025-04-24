// Package main is the main package for the updated the generated code.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"mime"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/conneroisu/conneroh.com/internal/cache"
	"github.com/conneroisu/conneroh.com/internal/llama"
	"github.com/conneroisu/conneroh.com/internal/logger"
	"github.com/conneroisu/conneroh.com/internal/markdown"
	"github.com/conneroisu/conneroh.com/internal/tigris"
	"github.com/rotisserie/eris"
	"github.com/spf13/afero"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	_ "modernc.org/sqlite"
)

var (
	workers = flag.Int("jobs", 40, "number of parallel workers")
)

func main() {
	flag.Parse()
	slog.SetDefault(logger.DefaultLogger)
	if err := run(
		context.Background(),
		os.Getenv,
		*workers,
	); err != nil {
		formattedStr := eris.ToString(err, true)
		fmt.Println(formattedStr)
		os.Exit(1)
	}
}

func run(
	ctx context.Context,
	getenv func(string) string,
	workers int,
) error {
	var (
		errCh = make(chan error)
		msgCh = make(cache.MsgChannel, workers)
		queCh = make(cache.MsgChannel, 5)
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
	err = assets.InitDB(ctx, db)
	if err != nil {
		return err
	}
	fs := afero.NewBasePathFs(afero.NewOsFs(), assets.VaultLoc)
	md := markdown.NewMD(fs)
	// db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	h := cache.NewHollywood(workers, fs, db, ti, ol, md)
	println("starting hollywood")

	go h.Start(innerCtx, msgCh, queCh, errCh, &wg)
	go ReadFS(fs, &wg, queCh, errCh)
	go Querer(innerCtx, queCh, msgCh)
	time.Sleep(time.Millisecond * 50)
	go func() {
		wg.Wait()
		slog.Info("hollywood finished")
		cancel()
	}()
	for {
		select {
		case <-time.After(time.Second * 1):
			slog.Info("waiting for hollywood to finish")
		case <-innerCtx.Done():
			slog.Info("inner context done")
			close(msgCh)
			close(queCh)
			close(errCh)
			return nil
		case err := <-errCh:
			slog.Error("error", "err", err)
			cancel()
			close(msgCh)
			close(queCh)
			close(errCh)
			return err
		}
	}
}

// ReadFS reads the filesystem and sends messages to the message channel.
func ReadFS(
	fs afero.Fs,
	wg *sync.WaitGroup,
	queCh cache.MsgChannel,
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
				slog.Info("skipping file", "path", fPath, "type", mime.TypeByExtension(filepath.Ext(fPath)))
				return nil
			}
			queCh <- cache.Msg{Path: fPath, Type: msgType}
			return nil
		})
	if err != nil {
		errCh <- fmt.Errorf("error walking filesystem: %w", err)
	}
}

// Querer reads from the message channel and sends messages to the que channel.
func Querer(
	ctx context.Context,
	queCh cache.MsgChannel,
	msgCh cache.MsgChannel,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-queCh:
			slog.Info("queing", "msg", msg)
			msgCh <- msg
		}
	}
}
