package cache

import (
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
	fs afero.Fs,
	wg waitGroup,
	queCh MsgChannel,
	errCh chan error,
) {
	wg.Add(1)
	defer wg.Done()
	err := afero.Walk(
		fs,
		".",
		func(fPath string, info os.FileInfo, err error) error {
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
			queCh <- Msg{Path: fPath, Type: msgType}
			return nil
		})
	if err != nil {
		errCh <- fmt.Errorf("error walking filesystem: %w", err)
	}
}
