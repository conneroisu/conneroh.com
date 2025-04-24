// Package main is the main package for the updated the generated code.
package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/conneroisu/conneroh.com/internal/cache"
	"github.com/conneroisu/conneroh.com/internal/cwg"
	"github.com/conneroisu/conneroh.com/internal/llama"
	"github.com/conneroisu/conneroh.com/internal/logger"
	"github.com/conneroisu/conneroh.com/internal/markdown"
	"github.com/conneroisu/conneroh.com/internal/tigris"
	"github.com/rotisserie/eris"
	"github.com/spf13/afero"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"
	_ "modernc.org/sqlite"
)

var (
	workers = flag.Int("jobs", 40, "number of parallel workers")
)

func main() {
	flag.Parse()
	slog.SetDefault(logger.DefaultProdLogger)

	// Create a context that will be canceled on interrupt signals
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP)
	defer stop()

	if err := Run(
		ctx,
		os.Getenv,
		*workers,
	); err != nil {
		formattedStr := eris.ToString(err, true)
		fmt.Println(formattedStr)
		os.Exit(1)
	}
}

// Run runs the main program, update.
func Run(
	ctx context.Context,
	getenv func(string) string,
	workers int,
) error {
	var (
		errCh = make(chan error)
		msgCh = make(cache.MsgChannel)
		queCh = make(cache.MsgChannel)
		// wg    = sync.WaitGroup{}
		wg cwg.DebugWaitGroup
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
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	h := cache.NewHollywood(workers, fs, db, ti, ol, md)
	go h.Start(innerCtx, msgCh, queCh, errCh, &wg)
	go cache.ReadFS(innerCtx, fs, &wg, queCh, errCh)
	go cache.Querer(innerCtx, queCh, msgCh)
	time.Sleep(time.Millisecond * 50)

	// Goroutine to handle wait group completion
	go func() {
		wg.Wait()
		slog.Debug("hollywood finished")
		cancel()
	}()

	// Main loop with graceful shutdown
	for {
		select {
		case err := <-errCh:
			if errors.Is(err, context.Canceled) {
				slog.Debug("received termination signal")
				cancel() // Cancel inner context
				return gracefulShutdown(msgCh, queCh, errCh)
			}
			slog.Error("error", "err", err)
			cancel()
			return gracefulShutdown(msgCh, queCh, errCh)

		case <-ctx.Done():
			slog.Info("received termination signal")
			cancel() // Cancel inner context
			return gracefulShutdown(msgCh, queCh, errCh)

		case <-innerCtx.Done():
			slog.Info("inner context done")
			return gracefulShutdown(msgCh, queCh, errCh)

		case <-time.After(time.Second * 1):
			slog.Info("waiting for hollywood to finish", slog.Int("count", wg.Count()))
			wg.PrintActiveDebugInfo()
		}
	}
}

// gracefulShutdown performs a clean shutdown by draining and closing channels
func gracefulShutdown(msgCh, queCh cache.MsgChannel, errCh chan error) error {
	slog.Info("performing graceful shutdown")

	// Drain error channel first to capture any final errors
	var errs error
drainLoop:
	for {
		select {
		case err := <-errCh:
			errs = eris.Wrapf(errs, "error while handling error: %v", err)
			slog.Error("error during shutdown", "err", err)
		case <-time.After(100 * time.Millisecond):
			break drainLoop
		}
	}

	// Close all channels
	close(msgCh)
	close(queCh)
	close(errCh)

	slog.Info("shutdown complete")
	return errs
}
