// Package cache provides a simplified caching system for assets
package cache

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

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

// TaskType represents different types of tasks
type TaskType string

const (
	// TypeAsset is a task for processing an asset file
	TypeAsset TaskType = "asset"
	// TypeDoc is a task for processing a document file
	TypeDoc TaskType = "doc"
	// TypeRelationship is a task for updating relationships
	TypeRelationship TaskType = "relationship"
)

// Task represents a unit of work to be processed
type Task struct {
	Type    TaskType
	Path    string
	Content []byte
	Doc     *assets.Doc
	// For relationship tasks
	RelationshipFn func(context.Context) error
	RetryCount     int
}

// Stats tracks statistics about task processing
type Stats struct {
	TotalSubmitted     atomic.Int64
	TotalCompleted     atomic.Int64
	TotalRelationships atomic.Int64
	lastActivityNs     atomic.Int64 // unixNano
	scanComplete       atomic.Bool
	shutdownInitiated  atomic.Bool
}

// RecordTaskSubmitted records that a task was submitted
func (s *Stats) RecordTaskSubmitted() {
	s.TotalSubmitted.Add(1)
	s.lastActivityNs.Store(time.Now().UnixNano())
}

// RecordTaskCompleted records that a task was completed
func (s *Stats) RecordTaskCompleted() {
	s.TotalCompleted.Add(1)
	s.lastActivityNs.Store(time.Now().UnixNano())
}

// RecordRelationshipTask records that a relationship task was submitted
func (s *Stats) RecordRelationshipTask() {
	s.TotalRelationships.Add(1)
	s.lastActivityNs.Store(time.Now().UnixNano())
}

// SetScanComplete marks that filesystem scanning is complete
func (s *Stats) SetScanComplete() {
	s.scanComplete.Store(true)
	slog.Debug("filesystem scan complete",
		"submitted", s.TotalSubmitted.Load(),
		"completed", s.TotalCompleted.Load(),
		"relationships", s.TotalRelationships.Load())
}

// IsComplete checks if all tasks are complete
func (s *Stats) IsComplete() bool {
	return s.scanComplete.Load() &&
		(s.TotalCompleted.Load() >= s.TotalSubmitted.Load()) &&
		!s.shutdownInitiated.Load()
}

// TimeSinceActivity returns the time since the last activity
func (s *Stats) TimeSinceActivity() time.Duration {
	lastNs := s.lastActivityNs.Load()
	return time.Since(time.Unix(0, lastNs))
}

// SetShutdownInitiated marks that shutdown has been initiated
func (s *Stats) SetShutdownInitiated() {
	s.shutdownInitiated.Store(true)
}

// LogStatus logs the current status
func (s *Stats) LogStatus() {
	lastNs := s.lastActivityNs.Load()
	slog.Debug("processor status",
		"submitted", s.TotalSubmitted.Load(),
		"completed", s.TotalCompleted.Load(),
		"relationships", s.TotalRelationships.Load(),
		"scanComplete", s.scanComplete.Load(),
		"idleTime", time.Since(time.Unix(0, lastNs)).String())
}

// Processor handles processing of assets and documents
type Processor struct {
	fs          afero.Fs
	db          *bun.DB
	s3Client    s3Client
	embedder    embedder
	md          goldmark.Markdown
	taskCh      chan Task
	errCh       chan error
	wg          *sync.WaitGroup
	memCache    *memCache
	entityCache *entityCache
	dbMutex     *sync.Mutex // Mutex for database operations
	batchOps    *batchOperations
	stats       *Stats
	doneCh      chan struct{} // Channel to signal completion
}

// Interface for S3 operations
type s3Client interface {
	PutObject(
		ctx context.Context,
		input *s3.PutObjectInput,
		opts ...func(*s3.Options),
	) (*s3.PutObjectOutput, error)
}

// Interface for embedding operations
type embedder interface {
	Embeddings(
		ctx context.Context,
		content string,
		doc *assets.Doc,
	) error
}

// batchOperations handles batched database operations
type batchOperations struct {
	mu             sync.Mutex
	cacheUpdates   []*assets.Cache
	batchSize      int
	db             *bun.DB
	updateInterval time.Duration
	shutdownCh     chan struct{}
}

// newBatchOperations creates a new batch operations manager
func newBatchOperations(db *bun.DB, batchSize int) *batchOperations {
	b := &batchOperations{
		cacheUpdates:   make([]*assets.Cache, 0, batchSize),
		batchSize:      batchSize,
		db:             db,
		updateInterval: 500 * time.Millisecond,
		shutdownCh:     make(chan struct{}),
	}

	// Start a single background goroutine for regular flushing
	go func() {
		ticker := time.NewTicker(b.updateInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				b.flush(context.Background())
			case <-b.shutdownCh:
				// Final flush on shutdown
				b.flush(context.Background())
				return
			}
		}
	}()

	return b
}

// scheduleCacheUpdate adds a cache update to the batch
func (b *batchOperations) scheduleCacheUpdate(ctx context.Context, cache *assets.Cache) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Add to queue (create a copy to avoid race conditions)
	cacheCopy := *cache
	b.cacheUpdates = append(b.cacheUpdates, &cacheCopy)

	// If we've reached batch size, flush immediately
	if len(b.cacheUpdates) >= b.batchSize {
		// Unlock while executing to avoid deadlock
		b.mu.Unlock()
		b.flush(ctx)
		b.mu.Lock()
	}
}

// flush immediately processes any pending updates regardless of batch size
func (b *batchOperations) flush(ctx context.Context) {
	b.mu.Lock()

	// Nothing to do
	if len(b.cacheUpdates) == 0 {
		b.mu.Unlock()
		return
	}

	// Take the current batch and reset
	updates := b.cacheUpdates
	b.cacheUpdates = make([]*assets.Cache, 0, b.batchSize)

	// Unlock while executing the batch
	b.mu.Unlock()

	slog.Debug("executing batch cache update", "count", len(updates))

	// Use a timeout context to prevent hanging
	execCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := b.db.NewInsert().
		Model(&updates).
		On("CONFLICT (path) DO UPDATE").
		Set("hashed = EXCLUDED.hashed").
		Set("x = EXCLUDED.x").
		Set("y = EXCLUDED.y").
		Set("z = EXCLUDED.z").
		Exec(execCtx)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("batch update timed out", "count", len(updates))
		} else {
			slog.Error("failed to batch update cache", "err", err, "count", len(updates))
		}
	}
}

// shutdown stops the background flush timer and performs a final flush
func (b *batchOperations) shutdown() {
	close(b.shutdownCh)
}

// memCache provides an in-memory cache to reduce database queries
type memCache struct {
	once  sync.Once
	mu    sync.RWMutex
	paths map[string]string // path -> hash
}

// newMemCache creates a new memory cache
func newMemCache() *memCache {
	return &memCache{
		paths: make(map[string]string),
	}
}

// Get gets a hash from the cache
func (m *memCache) Get(path string) (string, bool) {
	m.mu.RLock()
	hash, ok := m.paths[path]
	m.mu.RUnlock()
	return hash, ok
}

// Set sets a hash in the cache
func (m *memCache) Set(path, hash string) {
	m.mu.Lock()
	m.paths[path] = hash
	m.mu.Unlock()
}

// ensureLoaded ensures the cache is loaded from the database once
func (m *memCache) ensureLoaded(ctx context.Context, db *bun.DB) error {
	var loadErr error
	m.once.Do(func() {
		loadErr = m.loadFromDB(ctx, db)
	})
	return loadErr
}

// loadFromDB loads ALL cache entries from the database in a single query
func (m *memCache) loadFromDB(ctx context.Context, db *bun.DB) error {
	slog.Debug("loading all cache entries from database")

	// Use a single query to get all cache entries
	var caches []assets.Cache

	err := db.NewSelect().
		Model(&caches).
		Scan(ctx)

	if err != nil {
		return eris.Wrap(err, "failed to load cache from database")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Reset and populate the cache
	m.paths = make(map[string]string, len(caches))
	for _, cache := range caches {
		m.paths[cache.Path] = cache.Hash
	}

	slog.Debug("loaded all cache entries into memory", "count", len(caches))
	return nil
}

// entityCache provides a cache for entities to reduce database lookups
type entityCache struct {
	mu       sync.RWMutex
	tags     map[string]*assets.Tag
	posts    map[string]*assets.Post
	projects map[string]*assets.Project
}

// newEntityCache creates a new entity cache
func newEntityCache() *entityCache {
	return &entityCache{
		tags:     make(map[string]*assets.Tag),
		posts:    make(map[string]*assets.Post),
		projects: make(map[string]*assets.Project),
	}
}

// GetTag gets a tag from the cache
func (e *entityCache) GetTag(slug string) (*assets.Tag, bool) {
	e.mu.RLock()
	tag, ok := e.tags[slug]
	e.mu.RUnlock()
	return tag, ok
}

// SetTag sets a tag in the cache
func (e *entityCache) SetTag(tag *assets.Tag) {
	e.mu.Lock()
	tagCopy := *tag // Make a copy to avoid race conditions
	e.tags[tag.Slug] = &tagCopy
	e.mu.Unlock()
}

// GetPost gets a post from the cache
func (e *entityCache) GetPost(slug string) (*assets.Post, bool) {
	e.mu.RLock()
	post, ok := e.posts[slug]
	e.mu.RUnlock()
	return post, ok
}

// SetPost sets a post in the cache
func (e *entityCache) SetPost(post *assets.Post) {
	e.mu.Lock()
	postCopy := *post // Make a copy to avoid race conditions
	e.posts[post.Slug] = &postCopy
	e.mu.Unlock()
}

// GetProject gets a project from the cache
func (e *entityCache) GetProject(slug string) (*assets.Project, bool) {
	e.mu.RLock()
	project, ok := e.projects[slug]
	e.mu.RUnlock()
	return project, ok
}

// SetProject sets a project in the cache
func (e *entityCache) SetProject(project *assets.Project) {
	e.mu.Lock()
	projectCopy := *project // Make a copy to avoid race conditions
	e.projects[project.Slug] = &projectCopy
	e.mu.Unlock()
}

// Buffer pool for reusing buffers in markdown parsing
var bufferPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

// Static mapping of file extensions to content types to avoid reflection
var contentTypes = map[string]string{
	".jpg":   "image/jpeg",
	".jpeg":  "image/jpeg",
	".png":   "image/png",
	".gif":   "image/gif",
	".svg":   "image/svg+xml",
	".webp":  "image/webp",
	".css":   "text/css",
	".txt":   "text/plain",
	".md":    "text/markdown",
	".pdf":   "application/pdf",
	".xml":   "application/xml",
	".zip":   "application/zip",
	".mp3":   "audio/mpeg",
	".mp4":   "video/mp4",
	".webm":  "video/webm",
	".wav":   "audio/wav",
	".ico":   "image/x-icon",
	".woff":  "font/woff",
	".woff2": "font/woff2",
	".ttf":   "font/ttf",
	".otf":   "font/otf",
}

// getContentType returns the content type for a file extension
func getContentType(path string) string {
	ext := filepath.Ext(path)
	if contentType, ok := contentTypes[ext]; ok {
		return contentType
	}
	// Fallback to standard library for unknown extensions
	return "application/octet-stream"
}

// NewProcessor creates a new asset processor
func NewProcessor(
	fs afero.Fs,
	db *bun.DB,
	s3Client s3Client,
	embedder embedder,
	md goldmark.Markdown,
	bufferSize int,
) *Processor {
	memCache := newMemCache()
	entityCache := newEntityCache()

	// Initialize SQLite for better concurrency
	initDB(db)

	return &Processor{
		fs:          fs,
		db:          db,
		s3Client:    s3Client,
		embedder:    embedder,
		md:          md,
		taskCh:      make(chan Task, bufferSize),
		errCh:       make(chan error, bufferSize),
		wg:          &sync.WaitGroup{},
		memCache:    memCache,
		entityCache: entityCache,
		dbMutex:     &sync.Mutex{},
		batchOps:    newBatchOperations(db, 50), // Batch size of 50
		stats:       &Stats{lastActivityNs: atomic.Int64{}},
		doneCh:      make(chan struct{}),
	}
}

// initDB configures SQLite for better concurrency handling
func initDB(db *bun.DB) {
	// Set pragmas for better concurrent performance
	pragmas := []string{
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = NORMAL",
		"PRAGMA busy_timeout = 5000",
		"PRAGMA cache_size = 10000",
		"PRAGMA temp_store = MEMORY",   // Added to use memory for temp storage
		"PRAGMA mmap_size = 268435456", // Added to use memory mapping (256MB)
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			slog.Error("failed to set pragma", "pragma", pragma, "err", err)
		}
	}
}

// Start starts the processor workers
func (p *Processor) Start(ctx context.Context, numWorkers int) {
	// Initialize the last activity timestamp
	p.stats.lastActivityNs.Store(time.Now().UnixNano())

	// Start completion monitor
	go p.monitorCompletion(ctx)

	// Start status logger
	go p.logStatus(ctx)

	// Start worker goroutines
	for i := range numWorkers {
		go p.worker(ctx, i)
	}

	slog.Debug("processor started", "workers", numWorkers)
}

// monitorCompletion monitors for task completion
func (p *Processor) monitorCompletion(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Debug("completion monitor shutdown due to context cancellation")
			return

		case <-ticker.C:
			// Check if all tasks are complete
			if p.stats.IsComplete() {
				slog.Debug("all tasks complete - EXITING PROGRAM",
					"total", p.stats.TotalCompleted.Load(),
					"relationships", p.stats.TotalRelationships.Load())

				// Use direct exit since channel signaling may be unreliable
				os.Exit(0)
			}

			// Check for inactivity timeout (15 seconds - reduced from 30)
			if p.stats.scanComplete.Load() && p.stats.TimeSinceActivity() > 15*time.Second {
				slog.Warn("inactivity timeout detected - EXITING PROGRAM",
					"idleTime", p.stats.TimeSinceActivity().String())

				// Use direct exit since channel signaling may be unreliable
				os.Exit(0)
			}
		}
	}
}

// logStatus logs the processor status periodically
func (p *Processor) logStatus(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.stats.LogStatus()

			// Add extra diagnostics
			slog.Debug("task channel status",
				"channelCapacity", cap(p.taskCh),
				"currentQueueSize", len(p.taskCh))
		}
	}
}

// worker processes tasks from the task channel
func (p *Processor) worker(ctx context.Context, id int) {
	slog.Debug("starting worker", "id", id)
	defer slog.Debug("worker stopped", "id", id)

	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-p.taskCh:
			if !ok {
				return // Channel closed
			}

			p.wg.Add(1)
			func() {
				defer p.wg.Done()

				// Log task processing
				slog.Debug("processing task", "task", task)

				// Process the task
				if err := p.processTask(ctx, task); err != nil {
					select {
					case <-ctx.Done():
						return
					case p.errCh <- eris.Wrapf(err, "worker %d: error processing %s", id, task.Path):
					}
				}

				// Log task completion
				slog.Debug("task processed", "task", task)

				// Record task completion
				p.stats.RecordTaskCompleted()
			}()
		}
	}
}

// SubmitTask adds a task to the task channel
func (p *Processor) SubmitTask(task Task) {
	select {
	case p.taskCh <- task:
		// Task submitted successfully
		p.stats.RecordTaskSubmitted()
		if task.Type == TypeRelationship {
			p.stats.RecordRelationshipTask()
		}
	case <-time.After(5 * time.Second):
		// This shouldn't happen with a properly sized buffer
		p.errCh <- eris.Errorf("timeout submitting task: %s", task.Path)
	}
}

// ScanFS scans the filesystem and submits tasks
func (p *Processor) ScanFS() error {
	slog.Debug("scanning filesystem")

	// Create a buffered list to batch submit tasks
	var tasks []Task

	err := afero.Walk(p.fs, ".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil // Skip directories
		}

		// Determine task type based on file type
		var taskType TaskType
		switch {
		case assets.IsAllowedMediaType(path):
			taskType = TypeAsset
		case assets.IsAllowedDocumentType(path):
			taskType = TypeDoc
		default:
			// Skip files that aren't assets or documents
			return nil
		}

		// Add to batch
		tasks = append(tasks, Task{
			Type: taskType,
			Path: path,
		})

		// Submit in batches of 100
		if len(tasks) >= 100 {
			for _, task := range tasks {
				p.SubmitTask(task)
			}
			tasks = tasks[:0] // Clear batch
		}

		return nil
	})

	// Submit any remaining tasks
	for _, task := range tasks {
		p.SubmitTask(task)
	}

	// Mark scan as complete
	p.stats.SetScanComplete()

	return err
}

// WaitForCompletion waits for all tasks to complete
func (p *Processor) WaitForCompletion(timeout time.Duration) bool {
	slog.Debug("waiting for task completion", "timeout", timeout.String())

	// Create a timer for the timeout
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// Wait for done signal or timeout
	select {
	case <-p.doneCh:
		// All tasks completed
		slog.Info("received completion signal - all tasks completed")
		return true
	case <-timer.C:
		// Timeout occurred
		p.stats.LogStatus()
		slog.Warn(
			"timeout waiting for completion",
			"timeout",
			timeout.String(),
			"queueSize",
			len(p.taskCh),
			"pendingTasks",
			p.stats.TotalSubmitted.Load()-p.stats.TotalCompleted.Load(),
		)
		return false
	}
}

// WaitForAllTasks waits for the wait group to finish
// This is different from WaitForCompletion as it waits for the actual
// wait group to be empty, not just for the done signal
func (p *Processor) WaitForAllTasks() {
	p.wg.Wait()
}

// Close properly cleans up all resources
func (p *Processor) Close() {
	// Shutdown the batch operations background goroutine
	p.batchOps.shutdown()

	// Close all channels
	close(p.taskCh)
	close(p.errCh)
	close(p.doneCh)

	slog.Debug("processor resources cleaned up")
}

// Errors returns the error channel
func (p *Processor) Errors() <-chan error {
	return p.errCh
}

// Done returns the done channel
func (p *Processor) Done() <-chan struct{} {
	return p.doneCh
}

// processTask handles processing of a single task
func (p *Processor) processTask(ctx context.Context, task Task) error {
	switch task.Type {
	case TypeAsset:
		return p.processAsset(ctx, task)
	case TypeDoc:
		return p.processDocument(ctx, task)
	case TypeRelationship:
		// Handle retry for relationship tasks
		if task.RetryCount > 3 {
			return eris.Errorf("max retries exceeded for relationship task: %s", task.Path)
		}

		if err := task.RelationshipFn(ctx); err != nil {
			// Schedule retry with backoff
			task.RetryCount++
			delay := time.Duration(task.RetryCount*250) * time.Millisecond

			go func() {
				select {
				case <-ctx.Done():
					return
				case <-time.After(delay):
					p.SubmitTask(task)
				}
			}()
			return nil // Don't propagate error as we're retrying
		}
		return nil
	default:
		return eris.Errorf("unknown task type: %s", task.Type)
	}
}

// processAsset processes an asset file
func (p *Processor) processAsset(ctx context.Context, task Task) error {
	// Read file if content wasn't provided
	var content []byte
	var err error

	if task.Content != nil {
		content = task.Content
	} else {
		content, err = afero.ReadFile(p.fs, task.Path)
		if err != nil {
			return eris.Wrapf(err, "failed to read asset: %s", task.Path)
		}
	}

	// Calculate hash
	hash := assets.Hash(content)

	// Check if asset is already cached and unchanged
	if isCached, err := p.checkCache(ctx, task.Path, hash); err != nil {
		return err
	} else if isCached {
		return nil // Skip if already cached and unchanged
	}

	// Upload to S3
	if err := p.uploadToS3(ctx, task.Path, content); err != nil {
		return err
	}

	// Update cache
	return p.updateCache(ctx, task.Path, hash)
}

// processDocument processes a document file
func (p *Processor) processDocument(ctx context.Context, task Task) error {
	// Read file if content wasn't provided
	var content []byte
	var err error

	if task.Content != nil {
		content = task.Content
	} else {
		content, err = afero.ReadFile(p.fs, task.Path)
		if err != nil {
			return eris.Wrapf(err, "failed to read document: %s", task.Path)
		}
	}

	// Calculate hash
	hash := assets.Hash(content)

	// Check if document is already cached and unchanged
	if isCached, err := p.checkCache(ctx, task.Path, hash); err != nil {
		return err
	} else if isCached {
		return nil // Skip if already cached and unchanged
	}

	// Parse document
	doc := &assets.Doc{
		Path: task.Path,
		Hash: hash,
	}

	// Set default values
	if err := assets.Defaults(doc); err != nil {
		return eris.Wrapf(err, "failed to set defaults for document: %s", task.Path)
	}

	// Parse markdown and frontmatter
	if err := p.parseMarkdown(content, doc); err != nil {
		return err
	}

	// Set slug based on path
	doc.Slug = assets.Slugify(task.Path)

	// Validate document
	if err := assets.Validate(task.Path, doc); err != nil {
		return eris.Wrapf(err, "document validation failed: %s", task.Path)
	}

	// Save document to database
	if err := p.saveDocument(ctx, doc); err != nil {
		return err
	}

	// Update cache
	return p.updateCache(ctx, task.Path, hash)
}

// checkCache checks if a file is already cached and unchanged
// This version uses the memory cache with sync.Once initialization
func (p *Processor) checkCache(ctx context.Context, path, hash string) (bool, error) {
	// Make sure cache is loaded (only once)
	if err := p.memCache.ensureLoaded(ctx, p.db); err != nil {
		return false, err
	}

	// Check in memory cache
	if cachedHash, ok := p.memCache.Get(path); ok {
		if cachedHash == hash {
			return true, nil
		}
		slog.Debug("file has changed", "path", path)
		return false, nil
	}

	// Not in cache
	return false, nil
}

// updateCache updates the cache record for a file
func (p *Processor) updateCache(ctx context.Context, path, hash string) error {
	cache := &assets.Cache{
		Path: path,
		Hash: hash,
	}

	// Update memory cache immediately
	p.memCache.Set(path, hash)

	// Schedule for batch update
	p.batchOps.scheduleCacheUpdate(ctx, cache)

	return nil
}

// uploadToS3 uploads a file to S3
func (p *Processor) uploadToS3(ctx context.Context, path string, data []byte) error {
	// Use custom content type function instead of mime.TypeByExtension
	contentType := getContentType(path)

	// Use a timeout context to prevent hanging on S3 operations
	uploadCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := p.s3Client.PutObject(uploadCtx, &s3.PutObjectInput{
		Bucket:      aws.String("conneroh"),
		Key:         aws.String(path),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return eris.New("S3 upload timed out")
		}
		return eris.Wrapf(err, "failed to upload to S3: %s", path)
	}

	slog.Debug("uploaded to S3", "path", path)
	return nil
}

// parseMarkdown parses markdown content and extracts frontmatter
func (p *Processor) parseMarkdown(content []byte, doc *assets.Doc) error {
	// Get a buffer from the pool
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	// Create a new parser context
	pCtx := parser.NewContext()

	// Convert markdown
	err := p.md.Convert(content, buf, parser.WithContext(pCtx))
	if err != nil {
		return eris.Wrapf(err, "failed to convert markdown: %s", doc.Path)
	}

	// Extract frontmatter
	metadata := frontmatter.Get(pCtx)
	if metadata == nil {
		return eris.Errorf("frontmatter is nil for %s", doc.Path)
	}

	if err := metadata.Decode(doc); err != nil {
		return eris.Wrapf(err, "failed to decode frontmatter: %s", doc.Path)
	}

	doc.Content = buf.String()
	return nil
}

// saveDocument saves a document to the database
func (p *Processor) saveDocument(ctx context.Context, doc *assets.Doc) error {
	switch {
	case strings.HasPrefix(doc.Path, assets.PostsLoc) && strings.HasSuffix(doc.Path, ".md"):
		return p.savePost(ctx, doc)
	case strings.HasPrefix(doc.Path, assets.ProjectsLoc) && strings.HasSuffix(doc.Path, ".md"):
		return p.saveProject(ctx, doc)
	case strings.HasPrefix(doc.Path, assets.TagsLoc) && strings.HasSuffix(doc.Path, ".md"):
		return p.saveTag(ctx, doc)
	default:
		return eris.Errorf("unknown document type: %s", doc.Path)
	}
}

// savePost saves a post to the database
func (p *Processor) savePost(ctx context.Context, doc *assets.Doc) error {
	// Convert Doc to Post
	post := &assets.Post{}
	copygen.ToPost(post, doc)

	// Use database mutex to prevent concurrent writes
	p.dbMutex.Lock()
	defer p.dbMutex.Unlock()

	// Save post with a timeout context to prevent hanging
	insertCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Save post
	_, err := p.db.NewInsert().
		Model(post).
		On("CONFLICT (slug) DO UPDATE").
		Set("title = EXCLUDED.title").
		Set("description = EXCLUDED.description").
		Set("content = EXCLUDED.content").
		Set("banner_path = EXCLUDED.banner_path").
		Set("created_at = EXCLUDED.created_at").
		Set("x = EXCLUDED.x").
		Set("y = EXCLUDED.y").
		Set("z = EXCLUDED.z").
		Exec(insertCtx)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return eris.New("database operation timed out")
		}
		return eris.Wrapf(err, "failed to save post: %s", doc.Path)
	}

	// Schedule relationship update task with a delay to avoid database contention
	time.AfterFunc(100*time.Millisecond, func() {
		p.SubmitTask(Task{
			Type: TypeRelationship,
			Path: doc.Path,
			RelationshipFn: func(ctx context.Context) error {
				return p.updatePostRelationships(ctx, post)
			},
		})
	})

	return nil
}

// saveProject saves a project to the database
func (p *Processor) saveProject(ctx context.Context, doc *assets.Doc) error {
	// Convert Doc to Project
	project := &assets.Project{}
	copygen.ToProject(project, doc)

	// Use database mutex to prevent concurrent writes
	p.dbMutex.Lock()
	defer p.dbMutex.Unlock()

	// Save project with a timeout context to prevent hanging
	insertCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Save project
	_, err := p.db.NewInsert().
		Model(project).
		On("CONFLICT (slug) DO UPDATE").
		Set("title = EXCLUDED.title").
		Set("description = EXCLUDED.description").
		Set("content = EXCLUDED.content").
		Set("banner_path = EXCLUDED.banner_path").
		Set("created_at = EXCLUDED.created_at").
		Set("x = EXCLUDED.x").
		Set("y = EXCLUDED.y").
		Set("z = EXCLUDED.z").
		Exec(insertCtx)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return eris.New("database operation timed out")
		}
		return eris.Wrapf(err, "failed to save project: %s", doc.Path)
	}

	// Schedule relationship update task with a delay to avoid database contention
	time.AfterFunc(10*time.Millisecond, func() {
		p.SubmitTask(Task{
			Type: TypeRelationship,
			Path: doc.Path,
			RelationshipFn: func(ctx context.Context) error {
				return p.updateProjectRelationships(ctx, project)
			},
		})
	})

	return nil
}

// saveTag saves a tag to the database
func (p *Processor) saveTag(ctx context.Context, doc *assets.Doc) error {
	// Convert Doc to Tag
	tag := &assets.Tag{}
	copygen.ToTag(tag, doc)

	// Use database mutex to prevent concurrent writes
	p.dbMutex.Lock()
	defer p.dbMutex.Unlock()

	// Save tag with a timeout context to prevent hanging
	insertCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Save tag
	_, err := p.db.NewInsert().
		Model(tag).
		On("CONFLICT (slug) DO UPDATE").
		Set("title = EXCLUDED.title").
		Set("description = EXCLUDED.description").
		Set("content = EXCLUDED.content").
		Set("banner_path = EXCLUDED.banner_path").
		Set("created_at = EXCLUDED.created_at").
		Set("x = EXCLUDED.x").
		Set("y = EXCLUDED.y").
		Set("z = EXCLUDED.z").
		Exec(insertCtx)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return eris.New("database operation timed out")
		}
		return eris.Wrapf(err, "failed to save tag: %s", doc.Path)
	}

	// Schedule relationship update task with a delay to avoid database contention
	time.AfterFunc(10*time.Millisecond, func() {
		p.SubmitTask(Task{
			Type: TypeRelationship,
			Path: doc.Path,
			RelationshipFn: func(ctx context.Context) error {
				return p.updateTagRelationships(ctx, tag)
			},
		})
	})

	return nil
}

// findTagBySlug finds a tag by its slug
func (p *Processor) findTagBySlug(
	ctx context.Context,
	slug, origin string,
) (*assets.Tag, error) {
	// First check in-memory cache
	if tag, ok := p.entityCache.GetTag(slug); ok {
		return tag, nil
	}

	// Use mutex to prevent concurrent database queries
	p.dbMutex.Lock()
	defer p.dbMutex.Unlock()

	// Use timeout to prevent hanging
	queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var tag assets.Tag

	err := p.db.NewSelect().
		Model(&tag).
		Where("slug = ?", slug).
		Scan(queryCtx)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, eris.New("database operation timed out")
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, eris.Errorf(
				"tag not found: %s (referenced from %s)",
				slug,
				origin,
			)
		}
		return nil, eris.Wrapf(err, "failed to find tag: %s", slug)
	}

	// Update cache
	p.entityCache.SetTag(&tag)

	return &tag, nil
}

// findPostBySlug finds a post by its slug
func (p *Processor) findPostBySlug(ctx context.Context, slug, origin string) (*assets.Post, error) {
	// First check in-memory cache
	if post, ok := p.entityCache.GetPost(slug); ok {
		return post, nil
	}

	// Use mutex to prevent concurrent database queries
	p.dbMutex.Lock()
	defer p.dbMutex.Unlock()

	// Use timeout to prevent hanging
	queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var post assets.Post

	err := p.db.NewSelect().
		Model(&post).
		Where("slug = ?", slug).
		Scan(queryCtx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, eris.New("database operation timed out")
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, eris.Errorf(
				"post not found: %s (referenced from %s)",
				slug,
				origin,
			)
		}
		return nil, eris.Wrapf(err, "failed to find post: %s", slug)
	}

	// Update cache
	p.entityCache.SetPost(&post)

	return &post, nil
}

// findProjectBySlug finds a project by its slug
func (p *Processor) findProjectBySlug(ctx context.Context, slug, origin string) (*assets.Project, error) {
	// First check in-memory cache
	if project, ok := p.entityCache.GetProject(slug); ok {
		return project, nil
	}

	// Use mutex to prevent concurrent database queries
	p.dbMutex.Lock()
	defer p.dbMutex.Unlock()

	// Use timeout to prevent hanging
	queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var project assets.Project

	err := p.db.NewSelect().
		Model(&project).
		Where("slug = ?", slug).
		Scan(queryCtx)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, eris.New("database operation timed out")
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, eris.Errorf(
				"project not found: %s (referenced from %s)",
				slug,
				origin,
			)
		}
		return nil, eris.Wrapf(err, "failed to find project: %s", slug)
	}

	// Update cache
	p.entityCache.SetProject(&project)

	return &project, nil
}

// updatePostRelationships updates relationships for a post
func (p *Processor) updatePostRelationships(
	ctx context.Context,
	post *assets.Post,
) error {
	for _, tagSlug := range post.TagSlugs { // Create tag relationships
		tag, err := p.findTagBySlug(ctx, tagSlug, post.Slug)
		if err != nil {
			return err
		}

		_, err = p.db.NewInsert().
			Model(&assets.PostToTag{
				PostID: post.ID,
				TagID:  tag.ID,
			}).
			On("CONFLICT (post_id, tag_id) DO NOTHING").
			Exec(ctx)
		if err != nil {
			return eris.Wrapf(
				err,
				"failed to create post-tag relationship: %s -> %s",
				post.Slug,
				tagSlug,
			)
		}
	}

	for _, postSlug := range post.PostSlugs { // Create post relationships
		relatedPost, err := p.findPostBySlug(ctx, postSlug, post.Slug)
		if err != nil {
			return err
		}

		_, err = p.db.NewInsert().
			Model(&assets.PostToPost{
				SourcePostID: post.ID,
				TargetPostID: relatedPost.ID,
			}).
			On("CONFLICT (source_post_id, target_post_id) DO NOTHING").
			Exec(ctx)

		if err != nil {
			return eris.Wrapf(
				err,
				"failed to create post-post relationship: %s -> %s",
				post.Slug,
				postSlug,
			)
		}
	}

	for _, projectSlug := range post.ProjectSlugs { // Create project relationships
		project, err := p.findProjectBySlug(ctx, projectSlug, post.Slug)
		if err != nil {
			return err
		}

		_, err = p.db.NewInsert().
			Model(&assets.PostToProject{
				PostID:    post.ID,
				ProjectID: project.ID,
			}).
			On("CONFLICT (post_id, project_id) DO NOTHING").
			Exec(ctx)

		if err != nil {
			return eris.Wrapf(
				err,
				"failed to create post-project relationship: %s -> %s",
				post.Slug,
				projectSlug,
			)
		}
	}

	return nil
}

// updateProjectRelationships updates relationships for a project
func (p *Processor) updateProjectRelationships(
	ctx context.Context,
	project *assets.Project,
) error {
	for _, projectSlug := range project.ProjectSlugs { // Create project relationships
		relatedProject, err := p.findProjectBySlug(ctx, projectSlug, project.Slug)
		if err != nil {
			return err
		}

		_, err = p.db.NewInsert().
			Model(&assets.ProjectToProject{
				SourceProjectID: project.ID,
				TargetProjectID: relatedProject.ID,
			}).
			On("CONFLICT (source_project_id, target_project_id) DO NOTHING").
			Exec(ctx)
		if err != nil {
			return eris.Wrapf(
				err,
				"failed to create project-project relationship: %s -> %s",
				project.Slug,
				projectSlug,
			)
		}
	}

	for _, tagSlug := range project.TagSlugs { // Create tag relationships
		tag, err := p.findTagBySlug(ctx, tagSlug, project.Slug)
		if err != nil {
			return err
		}

		_, err = p.db.NewInsert().
			Model(&assets.ProjectToTag{
				ProjectID: project.ID,
				TagID:     tag.ID,
			}).
			On("CONFLICT (project_id, tag_id) DO NOTHING").
			Exec(ctx)
		if err != nil {
			return eris.Wrapf(
				err,
				"failed to create project-tag relationship: %s -> %s",
				project.Slug,
				tagSlug,
			)
		}
	}

	for _, postSlug := range project.PostSlugs { // Create post relationships
		post, err := p.findPostBySlug(ctx, postSlug, project.Slug)
		if err != nil {
			return err
		}

		_, err = p.db.NewInsert().
			Model(&assets.PostToProject{
				ProjectID: project.ID,
				PostID:    post.ID,
			}).
			On("CONFLICT (project_id, post_id) DO NOTHING").
			Exec(ctx)
		if err != nil {
			return eris.Wrapf(
				err,
				"failed to create project-post relationship: %s -> %s",
				project.Slug,
				postSlug,
			)
		}
	}

	return nil
}

// updateTagRelationships updates relationships for a tag
func (p *Processor) updateTagRelationships(
	ctx context.Context,
	tag *assets.Tag,
) error {
	for _, tagSlug := range tag.TagSlugs { // Create tag relationships
		relatedTag, err := p.findTagBySlug(ctx, tagSlug, tag.Slug)
		if err != nil {
			return err
		}
		_, err = p.db.NewInsert().
			Model(&assets.TagToTag{
				SourceTagID: tag.ID,
				TargetTagID: relatedTag.ID,
			}).
			On("CONFLICT (source_tag_id, target_tag_id) DO NOTHING").
			Exec(ctx)
		if err != nil {
			return eris.Wrapf(
				err,
				"failed to create tag-tag relationship: %s -> %s",
				tag.Slug,
				tagSlug,
			)
		}
	}

	for _, postSlug := range tag.PostSlugs { // Create post relationships
		post, err := p.findPostBySlug(ctx, postSlug, tag.Slug)
		if err != nil {
			return err
		}

		_, err = p.db.NewInsert().
			Model(&assets.PostToTag{
				TagID:  tag.ID,
				PostID: post.ID,
			}).
			On("CONFLICT (tag_id, post_id) DO NOTHING").
			Exec(ctx)
		if err != nil {
			return eris.Wrapf(
				err,
				"failed to create tag-post relationship: %s -> %s",
				tag.Slug,
				postSlug,
			)
		}
	}

	for _, projectSlug := range tag.ProjectSlugs { // Create project relationships
		project, err := p.findProjectBySlug(ctx, projectSlug, tag.Slug)
		if err != nil {
			return err
		}

		_, err = p.db.NewInsert().
			Model(&assets.ProjectToTag{
				TagID:     tag.ID,
				ProjectID: project.ID,
			}).
			On("CONFLICT (tag_id, project_id) DO NOTHING").
			Exec(ctx)
		if err != nil {
			return eris.Wrapf(
				err,
				"failed to create tag-project relationship: %s -> %s",
				tag.Slug,
				projectSlug,
			)
		}
	}

	return nil
}
