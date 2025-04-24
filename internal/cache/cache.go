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
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/conneroisu/conneroh.com/internal/copygen"
	"github.com/conneroisu/conneroh.com/internal/tigris"
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

	// OllamaClient is the ollama client.
	OllamaClient interface {
		Embeddings(
			ctx context.Context,
			content string,
			emb *assets.Doc,
		) (err error)
	}
)

// CtxSend sends a message to the channel with a context.
func CtxSend[T any](ctx context.Context, ch chan T, msg T, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	select {
	case <-ctx.Done():
		return
	case ch <- msg:
		return
	}
}

// NewHollywood creates a new asset hollywood.
func NewHollywood(
	num int,
	fs afero.Fs,
	db *bun.DB,
	ti tigris.Client,
	ol OllamaClient,
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
	msgCh chan Msg,
	errCh chan error,
	wg *sync.WaitGroup,
) {
	for _, actor := range h {
		go actor.Start(ctx, msgCh, errCh, wg)
	}
}

// Actor is the actor for the document collection.
type Actor struct {
	msgCh MsgChannel
	fs    afero.Fs
	db    *bun.DB
	ti    tigris.Client
	ol    OllamaClient
	md    goldmark.Markdown
}

// NewActor creates a new asset actor.
func NewActor(
	fs afero.Fs,
	db *bun.DB,
	ti tigris.Client,
	ol OllamaClient,
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
	errCh chan error,
	wg *sync.WaitGroup,
) {
	a.msgCh = msgCh
	for {
		select {
		case <-time.After(time.Second * 5):
			slog.Debug("cache actor", "killed", "was idle for 5 seconds")
		case <-ctx.Done():
			return
		case msg := <-msgCh:
			wg.Add(1)
			err := a.Handle(ctx, wg, msg)
			if err != nil {
				go CtxSend(ctx, errCh, err, wg)
			}
			wg.Done()
		}
	}
}

// Handle handles a message to the asset actor.
// When a message is received, it is expected that the path is a valid asset path.
// (e.g. /assets/foo.md, foo.svg, etc)
func (a *Actor) Handle(
	ctx context.Context,
	wg *sync.WaitGroup,
	msg Msg,
) (err error) {
	var (
		doc     assets.Doc
		ok      bool
		content []byte
	)
	if msg.Path == "" {
		return eris.New("path is empty")
	}
	if msg.Type == "" {
		return eris.New("type is empty")
	}
	switch msg.Type {
	case MsgTypeAction:
		slog.Info("action", "tries", msg.Tries)
		err = msg.fn()
		msg.Tries++
		if msg.Tries > 3 {
			return eris.Wrapf(err, "failed to handle (path: %s) action (tries: %d)", msg.Path, msg.Tries)
		}
		go CtxSend(ctx, a.msgCh, msg, wg)
		return nil
	case MsgTypeAsset:
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
		err = Upload(ctx, a.ti, msg.Path, content)
		if err != nil {
			return eris.Wrapf(err, "failed to upload asset: %s", msg.Path)
		}
		var asset assets.Asset
		copygen.ToAsset(&asset, &doc)
		err = SaveAsset(ctx, a.db, msg.Path, &asset, a.msgCh, wg)
		if err != nil {
			return eris.Wrapf(err, "failed to save asset: %s", msg.Path)
		}
		return nil
	case MsgTypeDoc:
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
		var metadata *frontmatter.Data
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

		// Set slug and content
		doc.Slug = assets.Slugify(msg.Path)
		doc.Content = buf.String()

		err = assets.Validate(msg.Path, &doc)
		if err != nil {
			return eris.Wrapf(err, "failed to validate %s", msg.Path)
		}

		err = Save(ctx, a.db, msg.Path, &doc, a.msgCh, wg)
		if err != nil {
			return eris.Wrapf(err, "failed to save %s", msg.Path)
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
	client tigris.Client,
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
		// Log error but let errgroup handle it
		return fmt.Errorf("failed to upload asset %s: %w", path, err)
	}

	slog.Info("asset upload successful", "path", path)
	return nil
}

func (a *Actor) isCached(ctx context.Context, msg Msg, doc *assets.Doc) (bool, error) {
	// if content has not changed, skip
	switch {
	case strings.HasPrefix(msg.Path, assets.PostsLoc):
		var p assets.Post
		_, err := a.db.NewSelect().
			Model(assets.EmpPost).
			Where("path = ?", msg.Path).
			Exec(ctx, &p)
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		if p.Hash == doc.Hash {
			return true, nil
		}
	case strings.HasPrefix(msg.Path, assets.ProjectsLoc):
		var p assets.Project
		_, err := a.db.NewSelect().
			Model(assets.EmpProject).
			Where("path = ?", msg.Path).
			Exec(ctx, &p)
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		if p.Hash == doc.Hash {
			return true, nil
		}
	case strings.HasPrefix(msg.Path, assets.TagsLoc):
		var t assets.Tag
		_, err := a.db.NewSelect().
			Model(assets.EmpTag).
			Where("path = ?", msg.Path).
			Exec(ctx, &t)
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		if t.Hash == doc.Hash {
			return true, nil
		}
	case strings.HasPrefix(msg.Path, assets.AssetsLoc):
		var ass assets.Asset
		_, err := a.db.NewSelect().
			Model(assets.EmpAsset).
			Where("path = ?", msg.Path).
			Exec(ctx, &a)
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		if ass.Hash == doc.Hash {
			return true, nil
		}
	}
	return false, nil
}
