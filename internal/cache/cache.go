package cache

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"mime"
	"path/filepath"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/conneroisu/conneroh.com/internal/data/db"
	"github.com/conneroisu/conneroh.com/internal/tigris"
	"github.com/rotisserie/eris"
	"github.com/spf13/afero"
	"github.com/uptrace/bun"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

// MsgChannel is a cache message channel.
type MsgChannel chan *Msg

// Hollywood is a slice of actors.
type Hollywood []*Actor

// OllamaClient is the ollama client.
type OllamaClient interface {
	Embeddings(
		ctx context.Context,
		content string,
		emb *assets.Doc,
	) (err error)
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
	msgCh chan *Msg,
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
)

// Msg is the message for the document actor.
type Msg struct {
	Type  MsgType
	Path  string
	Doc   *assets.Doc
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
		case <-ctx.Done():
			return
		case msg := <-msgCh:
			wg.Add(1)
			err := a.Handle(ctx, msg)
			if err != nil {
				errCh <- err
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
	msg *Msg,
) (err error) {
	var doc assets.Doc
	if msg.Path == "" {
		return fmt.Errorf("path is empty")
	}
	doc.RawContent, err = afero.ReadFile(a.fs, msg.Path)
	if err != nil {
		return err
	}
	switch msg.Type {
	case MsgTypeAsset:
		return Upload(ctx, a.ti, msg.Path, doc.RawContent)
	case MsgTypeDoc:
		var metadata *frontmatter.Data
		pCtx := parser.NewContext()
		buf := bytes.NewBufferString("")

		err = assets.Defaults(&doc)
		if err != nil {
			return err
		}

		err = a.md.Convert(doc.RawContent, buf, parser.WithContext(pCtx))
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
		doc.Hash = assets.Hash(doc.RawContent)

		err = assets.Validate(msg.Path, &doc)
		if err != nil {
			return eris.Wrapf(err, "failed to validate %s", msg.Path)
		}

		err = db.Save(ctx, a.db, msg.Path, &doc)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown message type: %s", msg.Type)
	}
	return nil
}

// UploadManager manages concurrent uploads using multiple mutexes
type UploadManager struct {
	mutexes    []sync.Mutex
	mutexCount int
	inUse      []bool
	lock       sync.Mutex // Meta-mutex to protect the inUse array
}

// NewUploadManager creates a new UploadManager with the specified number of mutexes
func NewUploadManager(count int) *UploadManager {
	if count <= 0 {
		count = 10 // Default to 10 mutexes if invalid count provided
	}

	manager := &UploadManager{
		mutexes:    make([]sync.Mutex, count),
		mutexCount: count,
		inUse:      make([]bool, count),
	}

	return manager
}

// AcquireMutex finds an available mutex, locks it, and returns its index
// This function will block until a mutex is available
func (um *UploadManager) AcquireMutex() int {
	for {
		um.lock.Lock()
		for i := range um.mutexCount {
			if !um.inUse[i] {
				// Mark this mutex as in use
				um.inUse[i] = true
				um.lock.Unlock()

				// Lock the mutex
				um.mutexes[i].Lock()
				return i
			}
		}
		// No mutex available, release meta-lock and wait a bit before trying again
		um.lock.Unlock()
		time.Sleep(10 * time.Millisecond)
	}
}

// ReleaseMutex unlocks the mutex at the specified index and marks it as available
func (um *UploadManager) ReleaseMutex(index int) {
	if index < 0 || index >= um.mutexCount {
		return
	}

	// First unlock the actual mutex
	um.mutexes[index].Unlock()

	// Then mark it as not in use
	um.lock.Lock()
	um.inUse[index] = false
	um.lock.Unlock()
}

// Global upload manager instance
var uploadManager = NewUploadManager(20) // Create 20 mutexes for more concurrency

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

	// Acquire a mutex for this upload - this will block until one is available
	mutexIndex := uploadManager.AcquireMutex()

	// Ensure we release the mutex when done
	defer uploadManager.ReleaseMutex(mutexIndex)

	// Check if context is done before proceeding with upload
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Continue with upload
	}

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

	slog.Info("asset upload successful", "path", path, "mutexIndex", mutexIndex)
	return nil
}
