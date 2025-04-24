package cache

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"mime"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/conneroisu/conneroh.com/internal/copygen"
	"github.com/rotisserie/eris"
	"github.com/spf13/afero"
	"github.com/uptrace/bun"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

type (
	// MsgChannel is a cache message channel.
	MsgChannel chan Msg

	// Hollywood is a slice of actors.
	Hollywood []*Actor
)

// NewHollywood creates a new asset hollywood.
func NewHollywood(
	num int,
	fs afero.Fs,
	db *bun.DB,
	ti tigrisClient,
	ol ollamaClient,
	md goldmark.Markdown,
) Hollywood {
	var actors []*Actor
	for range num {
		actors = append(actors, NewActor(fs, db, ti, ol, md))
	}
	return Hollywood(actors)
}

// Start starts the asset hollywood.
func (h Hollywood) Start(
	ctx context.Context,
	msgCh MsgChannel,
	queCh MsgChannel,
	errCh chan error,
	wg waitGroup,
) {
	for _, actor := range h {
		go actor.Start(ctx, msgCh, queCh, errCh, wg)
	}
}

// Actor is the actor for the document collection.
type Actor struct {
	queCh MsgChannel
	fs    afero.Fs
	db    *bun.DB
	ti    tigrisClient
	ol    ollamaClient
	md    goldmark.Markdown
}

// NewActor creates a new asset actor.
func NewActor(
	fs afero.Fs,
	db *bun.DB,
	ti tigrisClient,
	ol ollamaClient,
	md goldmark.Markdown,
) *Actor {
	return &Actor{fs: fs, db: db, ti: ti, ol: ol, md: md}
}

// MsgType is the type of message for the document actor message.
type MsgType string

const (
	// MsgTypeAsset is the message type for an asset.
	MsgTypeAsset MsgType = "asset"
	// MsgTypeDoc is the message type for a document.
	MsgTypeDoc MsgType = "doc"
	// MsgTypeAction is the message type for an action.
	MsgTypeAction MsgType = "action"
)

// Msg is the message for the document actor.
type Msg struct {
	Type  MsgType
	Path  string
	fn    func() error
	Tries int
}

// Start starts the asset actor.
func (a *Actor) Start(
	ctx context.Context,
	msgCh MsgChannel,
	queCh MsgChannel,
	errCh chan error,
	wg waitGroup,
) {
	a.queCh = queCh
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-msgCh:
			if !ok {
				// Channel closed, exit gracefully
				slog.Info("Channel closed, exiting")
				return
			}
			func() {
				wg.Add(1)
				defer wg.Done()
				err := a.Handle(ctx, msg)
				if err != nil {
					err = CtxSend(ctx, errCh, err)
					if err != nil {
						slog.Error("failed to send error to error channel", "err", err)
						return
					}
				}
			}()
		}
	}
}

// Handle handles a message to the asset actor.
// When a message is received, it is expected that the path is a valid asset path.
// (e.g. /assets/foo.md, foo.svg, etc)
func (a *Actor) Handle(
	ctx context.Context,
	msg Msg,
) (err error) {
	var (
		doc assets.Doc
	)
	if msg.Path == "" {
		return eris.New("path is empty")
	}
	if msg.Type == "" {
		return eris.New("type is empty")
	}
	doc.Path = msg.Path
	switch msg.Type {
	case MsgTypeAction:
		slog.Info("action", "tries", msg.Tries)
		err = msg.fn()
		msg.Tries++
		if msg.Tries > 3 && err != nil {
			return eris.Wrapf(
				err,
				"failed to handle (path: %s) action (tries: %d)",
				msg.Path,
				msg.Tries,
			)
		}

		if err != nil {
			// Only requeue if there was an error and we haven't exceeded retry count
			select {
			case <-ctx.Done():
				if errors.Is(ctx.Err(), context.Canceled) {
					return nil
				}
				return ctx.Err()
			default:
				err = CtxSend(ctx, a.queCh, msg)
				if err != nil {
					return eris.Wrapf(err, "failed to send message to queue channel: %s", msg.Path)
				}
			}
		}
		return nil
	case MsgTypeAsset:
		var (
			content []byte
			cache   assets.Cache
			ok      bool
		)
		content, err = afero.ReadFile(a.fs, msg.Path)
		if err != nil {
			return err
		}
		doc.Hash = assets.Hash(content)
		ok, err = a.isCached(ctx, msg, &doc)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
		copygen.ToCache(&cache, &doc)
		err = Cache(ctx, a.db, &cache)
		if err != nil {
			return eris.Wrapf(err, "failed to cache %s", msg.Path)
		}
		err = Upload(ctx, a.ti, msg.Path, content)
		if err != nil {
			return eris.Wrapf(err, "failed to upload asset: %s", msg.Path)
		}
		return nil
	case MsgTypeDoc:
		var (
			metadata *frontmatter.Data
			content  []byte
			cache    assets.Cache
			ok       bool
		)
		content, err = afero.ReadFile(a.fs, msg.Path)
		if err != nil {
			return err
		}
		doc.Hash = assets.Hash(content)
		ok, err = a.isCached(ctx, msg, &doc)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
		pCtx := parser.NewContext()
		buf := bytes.NewBufferString("")

		err = assets.Defaults(&doc)
		if err != nil {
			return err
		}

		err = a.md.Convert(content, buf, parser.WithContext(pCtx))
		if err != nil {
			return eris.Wrapf(
				err,
				"failed to convert %s's markdown to HTML",
				msg.Path,
			)
		}

		metadata = frontmatter.Get(pCtx)
		if metadata == nil {
			return eris.Errorf(
				"frontmatter is nil for %s",
				msg.Path,
			)
		}

		err = metadata.Decode(&doc)
		if err != nil {
			return eris.Wrapf(
				err,
				"failed to decode frontmatter of %s",
				msg.Path,
			)
		}

		doc.Slug = assets.Slugify(msg.Path)
		doc.Content = buf.String()

		err = assets.Validate(msg.Path, &doc)
		if err != nil {
			return eris.Wrapf(err, "failed to validate %s", msg.Path)
		}

		err = Save(ctx, a.db, msg.Path, &doc, a.queCh)
		if err != nil {
			return eris.Wrapf(err, "failed to save %s", msg.Path)
		}

		copygen.ToCache(&cache, &doc)
		err = Cache(ctx, a.db, &cache)
		if err != nil {
			return eris.Wrapf(err, "failed to cache %s", msg.Path)
		}
	default:
		return eris.New("unknown message type: " + string(msg.Type))
	}
	return nil
}

// Upload uploads the provided asset to the specified bucket.
// This version supports concurrent uploads using multiple mutexes
func Upload(
	ctx context.Context,
	client tigrisClient,
	path string,
	data []byte,
) error {
	slog.Debug(
		"asset changed uploading...",
		"path", path,
	)
	extension := filepath.Ext(path)
	if extension == "" {
		return fmt.Errorf("failed to get extension for %s", path)
	}
	contentType := mime.TypeByExtension(extension)

	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String("conneroh"),
		Key:         aws.String(path),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return fmt.Errorf("failed to upload asset %s: %w", path, err)
	}

	slog.Info("asset upload successful", "path", path)
	return nil
}

// isCached checks if a document or asset is already cached.
// Returns:
// - cacheID > 0: Document exists and is cached (unchanged)
// - cacheID == 0 && err == nil: Document doesn't exist or has changed (needs update)
// - err != nil: An error occurred
func (a *Actor) isCached(ctx context.Context, msg Msg, doc *assets.Doc) (bool, error) {
	var c assets.Cache
	err := a.db.NewSelect().
		Model(&c).
		Where("path = ?", msg.Path).
		Scan(ctx)

	// Document not found in cache, need to add it
	if errors.Is(err, sql.ErrNoRows) {
		slog.Info("asset not found in cache", "path", msg.Path)
		return false, nil
	}

	// Other database error
	if err != nil {
		return false, fmt.Errorf("failed to check cache: %w", err)
	}

	// Document found and hash matches - no changes needed
	if c.Hash == doc.Hash {
		slog.Debug("asset already cached and unchanged", "path", msg.Path)
		return true, nil
	}

	// Document found but hash is different - needs update
	slog.Info("asset has changed", "path", msg.Path, "old_hash", c.Hash, "new_hash", doc.Hash)
	return false, nil
}

// Cache caches a document to the database.
func Cache(
	ctx context.Context,
	db *bun.DB,
	cache *assets.Cache,
) error {
	_, err := db.NewInsert().
		Model(cache).
		On("CONFLICT (path) DO UPDATE").
		Set("hashed = EXCLUDED.hashed").
		Set("x = EXCLUDED.x").
		Set("y = EXCLUDED.y").
		Set("z = EXCLUDED.z").
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to cache %s: %w", cache.Path, err)
	}
	return nil
}
