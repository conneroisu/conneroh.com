package cache

import (
	"context"
	"fmt"
	"log/slog"
	"mime"
	"os"
	"path/filepath"

	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/spf13/afero"
)

// ReadFS reads the filesystem and sends messages to the message channel.
func ReadFS(
	ctx context.Context, // Add context parameter
	fs afero.Fs,
	wg waitGroup,
	queCh MsgChannel,
	errCh chan error,
) {
	wg.Add(1)
	defer wg.Done()

	walkDone := make(chan struct{})
	walkErr := make(chan error, 1)

	// Start file walk in separate goroutine
	go func() {
		err := afero.Walk(
			fs,
			".",
			func(fPath string, info os.FileInfo, err error) error {
				// Check context before processing each file
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					// Continue processing
				}

				wg.Add(1)
				defer wg.Done()

				var msgType MsgType
				switch {
				case err != nil:
					return err
				case assets.IsAllowedMediaType(fPath):
					msgType = MsgTypeAsset
				case assets.IsAllowedDocumentType(fPath):
					msgType = MsgTypeDoc
				case info.IsDir():
					return nil
				default:
					slog.Info("skipping file", "path", fPath, "type", mime.TypeByExtension(filepath.Ext(fPath)))
					return nil
				}

				// Send with context awareness
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					return CtxSend(ctx, queCh, Msg{Path: fPath, Type: msgType})
				}
			})

		if err != nil {
			walkErr <- err
		}
		close(walkDone)
	}()

	// Wait for walk to finish or context to be canceled
	select {
	case <-ctx.Done():
		slog.Info("ReadFS: context canceled, exiting")
		return
	case err := <-walkErr:
		select {
		case <-ctx.Done():
			// Context already canceled, just exit
			return
		default:
			// Try to send error
			CtxSend(ctx, errCh, fmt.Errorf("error walking filesystem: %w", err))
		}
	case <-walkDone:
		slog.Info("ReadFS: walk completed successfully")
	}
}
