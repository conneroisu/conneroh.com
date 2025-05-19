// Package cache provides a simplified caching system for assets
package cache

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"os"
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

// TaskType represents different types of tasks.
type TaskType string

const (
	// TypeAsset is a task for processing an asset file.
	TypeAsset TaskType = "asset"
	// TypeDoc is a task for processing a document file.
	TypeDoc TaskType = "doc"
	// TypeRelationship is a task for updating relationships.
	TypeRelationship TaskType = "relationship"
)

// Task represents a unit of work to be processed.
type Task struct {
	Type    TaskType
	Path    string
	Content []byte
	Doc     *assets.Doc
	// For relationship tasks
	RelationshipFn func(context.Context) error
	RetryCount     int
}

// DBTaskType represents different types of database tasks.
type DBTaskType string

const (
	// DBSavePost is for saving a post.
	DBSavePost DBTaskType = "save_post"
	// DBSaveProject is for saving a project.
	DBSaveProject DBTaskType = "save_project"
	// DBSaveTag is for saving a tag.
	DBSaveTag DBTaskType = "save_tag"
	// DBFindTag is for finding a tag.
	DBFindTag DBTaskType = "find_tag"
	// DBFindPost is for finding a post.
	DBFindPost DBTaskType = "find_post"
	// DBFindProject is for finding a project.
	DBFindProject DBTaskType = "find_project"
	// DBUpdateCache is for updating the cache.
	DBUpdateCache DBTaskType = "update_cache"
	// DBLoadCache is for loading the cache.
	DBLoadCache DBTaskType = "load_cache"
	// DBUpdateRelationship is for updating relationships.
	DBUpdateRelationship DBTaskType = "update_relationship"
)

// DBTask represents a database operation to be processed by the dedicated DB goroutine.
type DBTask struct {
	Type         DBTaskType
	Data         interface{}
	ResponseChan chan interface{}
	ErrorChan    chan error
}

// Stats tracks statistics about task processing.
type Stats struct {
	TotalSubmitted     atomic.Int64
	TotalCompleted     atomic.Int64
	TotalRelationships atomic.Int64
	DBTasksSubmitted   atomic.Int64
	DBTasksCompleted   atomic.Int64
	lastActivityNs     atomic.Int64 // unixNano
	scanComplete       atomic.Bool
	shutdownInitiated  atomic.Bool
}

// RecordTaskSubmitted records that a task was submitted.
func (s *Stats) RecordTaskSubmitted() {
	s.TotalSubmitted.Add(1)
	s.lastActivityNs.Store(time.Now().UnixNano())
}

// RecordTaskCompleted records that a task was completed.
func (s *Stats) RecordTaskCompleted() {
	s.TotalCompleted.Add(1)
	s.lastActivityNs.Store(time.Now().UnixNano())
}

// RecordDBTaskSubmitted records that a DB task was submitted.
func (s *Stats) RecordDBTaskSubmitted() {
	s.DBTasksSubmitted.Add(1)
	s.lastActivityNs.Store(time.Now().UnixNano())
}

// RecordDBTaskCompleted records that a DB task was completed.
func (s *Stats) RecordDBTaskCompleted() {
	s.DBTasksCompleted.Add(1)
	s.lastActivityNs.Store(time.Now().UnixNano())
}

// RecordRelationshipTask records that a relationship task was submitted.
func (s *Stats) RecordRelationshipTask() {
	s.TotalRelationships.Add(1)
	s.lastActivityNs.Store(time.Now().UnixNano())
}

// SetScanComplete marks that filesystem scanning is complete.
func (s *Stats) SetScanComplete() {
	s.scanComplete.Store(true)
	slog.Debug("filesystem scan complete",
		"submitted", s.TotalSubmitted.Load(),
		"completed", s.TotalCompleted.Load(),
		"relationships", s.TotalRelationships.Load(),
		"dbTasksSubmitted", s.DBTasksSubmitted.Load(),
		"dbTasksCompleted", s.DBTasksCompleted.Load())
}

// IsComplete checks if all tasks are complete.
func (s *Stats) IsComplete() bool {
	return s.scanComplete.Load() &&
		(s.TotalCompleted.Load() >= s.TotalSubmitted.Load()) &&
		(s.DBTasksCompleted.Load() >= s.DBTasksSubmitted.Load()) &&
		!s.shutdownInitiated.Load()
}

// TimeSinceActivity returns the time since the last activity.
func (s *Stats) TimeSinceActivity() time.Duration {
	lastNs := s.lastActivityNs.Load()

	return time.Since(time.Unix(0, lastNs))
}

// SetShutdownInitiated marks that shutdown has been initiated.
func (s *Stats) SetShutdownInitiated() {
	s.shutdownInitiated.Store(true)
}

// LogStatus logs the current status.
func (s *Stats) LogStatus() {
	lastNs := s.lastActivityNs.Load()
	slog.Debug("processor status",
		"submitted", s.TotalSubmitted.Load(),
		"completed", s.TotalCompleted.Load(),
		"relationships", s.TotalRelationships.Load(),
		"dbTasksSubmitted", s.DBTasksSubmitted.Load(),
		"dbTasksCompleted", s.DBTasksCompleted.Load(),
		"scanComplete", s.scanComplete.Load(),
		"idleTime", time.Since(time.Unix(0, lastNs)).String())
}

// Processor handles processing of assets and documents.
type Processor struct {
	fs           afero.Fs
	db           *bun.DB
	s3Client     s3Client
	s3Bucket     string
	embedder     embedder
	md           goldmark.Markdown
	taskCh       chan Task
	dbTaskCh     chan DBTask
	errCh        chan error
	wg           *sync.WaitGroup
	memCache     *memCache
	entityCache  *entityCache
	stats        *Stats
	doneCh       chan struct{} // Channel to signal completion
	dbShutdownCh chan struct{} // Channel to signal DB goroutine shutdown
}

// Interface for S3 operations.
type s3Client interface {
	PutObject(
		ctx context.Context,
		input *s3.PutObjectInput,
		opts ...func(*s3.Options),
	) (*s3.PutObjectOutput, error)
}

// Interface for embedding operations.
type embedder interface {
	Embeddings(
		ctx context.Context,
		content string,
		doc *assets.Doc,
	) error
}

// memCache provides an in-memory cache to reduce database queries.
type memCache struct {
	mu    sync.RWMutex
	paths map[string]string // path -> hash
}

// newMemCache creates a new memory cache.
func newMemCache() *memCache {
	return &memCache{
		paths: make(map[string]string),
	}
}

// Get gets a hash from the cache.
func (m *memCache) Get(path string) (string, bool) {
	m.mu.RLock()
	hash, ok := m.paths[path]
	m.mu.RUnlock()

	return hash, ok
}

// Set sets a hash in the cache.
func (m *memCache) Set(path, hash string) {
	m.mu.Lock()
	m.paths[path] = hash
	m.mu.Unlock()
}

// loadFromDB loads ALL cache entries from the database in a single query.
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

// entityCache provides a cache for entities to reduce database lookups.
type entityCache struct {
	mu       sync.RWMutex
	tags     map[string]*assets.Tag
	posts    map[string]*assets.Post
	projects map[string]*assets.Project
}

// newEntityCache creates a new entity cache.
func newEntityCache() *entityCache {
	return &entityCache{
		tags:     make(map[string]*assets.Tag),
		posts:    make(map[string]*assets.Post),
		projects: make(map[string]*assets.Project),
	}
}

// GetTag gets a tag from the cache.
func (e *entityCache) GetTag(slug string) (*assets.Tag, bool) {
	e.mu.RLock()
	tag, ok := e.tags[slug]
	e.mu.RUnlock()

	return tag, ok
}

// SetTag sets a tag in the cache.
func (e *entityCache) SetTag(tag *assets.Tag) {
	e.mu.Lock()
	tagCopy := *tag // Make a copy to avoid race conditions
	e.tags[tag.Slug] = &tagCopy
	e.mu.Unlock()
}

// GetPost gets a post from the cache.
func (e *entityCache) GetPost(slug string) (*assets.Post, bool) {
	e.mu.RLock()
	post, ok := e.posts[slug]
	e.mu.RUnlock()

	return post, ok
}

// SetPost sets a post in the cache.
func (e *entityCache) SetPost(post *assets.Post) {
	e.mu.Lock()
	postCopy := *post // Make a copy to avoid race conditions
	e.posts[post.Slug] = &postCopy
	e.mu.Unlock()
}

// GetProject gets a project from the cache.
func (e *entityCache) GetProject(slug string) (*assets.Project, bool) {
	e.mu.RLock()
	project, ok := e.projects[slug]
	e.mu.RUnlock()

	return project, ok
}

// SetProject sets a project in the cache.
func (e *entityCache) SetProject(project *assets.Project) {
	e.mu.Lock()
	projectCopy := *project // Make a copy to avoid race conditions
	e.projects[project.Slug] = &projectCopy
	e.mu.Unlock()
}

// Buffer pool for reusing buffers in markdown parsing.
var bufferPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

// NewProcessor creates a new asset processor.
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

	bucketName := os.Getenv("BUCKET_NAME")
	if bucketName == "" {
		panic("BUCKET_NAME environment variable is not set")
	}

	return &Processor{
		fs:           fs,
		db:           db,
		s3Client:     s3Client,
		s3Bucket:     bucketName,
		embedder:     embedder,
		md:           md,
		taskCh:       make(chan Task, bufferSize),
		dbTaskCh:     make(chan DBTask, bufferSize), // DB task channel with buffer
		errCh:        make(chan error, bufferSize),
		wg:           &sync.WaitGroup{},
		memCache:     memCache,
		entityCache:  entityCache,
		stats:        &Stats{lastActivityNs: atomic.Int64{}},
		doneCh:       make(chan struct{}),
		dbShutdownCh: make(chan struct{}),
	}
}

// Start starts the processor workers.
func (p *Processor) Start(ctx context.Context, numWorkers int) {
	// Initialize the last activity timestamp
	p.stats.lastActivityNs.Store(time.Now().UnixNano())

	// Start completion monitor
	go p.monitorCompletion(ctx)

	// Start status logger
	go p.logStatus(ctx)

	// Start the dedicated database worker
	go p.dbWorker(ctx)

	// Start worker goroutines
	for i := range numWorkers {
		go p.worker(ctx, i)
	}

	slog.Debug("processor started", "workers", numWorkers)
}

// monitorCompletion monitors for task completion.
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
					"dbTasks", p.stats.DBTasksCompleted.Load(),
					"relationships", p.stats.TotalRelationships.Load())

				// Signal completion instead of direct exit
				p.stats.SetShutdownInitiated()
				close(p.doneCh)
			}

			// Check for inactivity timeout (15 seconds - reduced from 30)
			if p.stats.scanComplete.Load() && p.stats.TimeSinceActivity() > 15*time.Second {
				slog.Warn("inactivity timeout detected - EXITING PROGRAM",
					"idleTime", p.stats.TimeSinceActivity().String())

				// Signal completion instead of direct exit
				p.stats.SetShutdownInitiated()
				close(p.doneCh)
			}
		}
	}
}

// logStatus logs the processor status periodically.
func (p *Processor) logStatus(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.stats.LogStatus()
			slog.Debug("task channel status",
				"taskChannelCapacity", cap(p.taskCh),
				"currentTaskQueueSize", len(p.taskCh),
				"dbChannelCapacity", cap(p.dbTaskCh),
				"currentDBQueueSize", len(p.dbTaskCh))
		}
	}
}

// flushCacheBatch flushes the batch of cache updates to the database.
func (p *Processor) flushCacheBatch(ctx context.Context, batch []*assets.Cache) []*assets.Cache {
	if len(batch) == 0 {
		return batch
	}

	slog.Debug("executing batch cache update", "count", len(batch))

	// Use a timeout context to prevent hanging
	execCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := p.db.NewInsert().
		Model(&batch).
		On("CONFLICT (path) DO UPDATE").
		Set("hashed = EXCLUDED.hashed").
		Set("x = EXCLUDED.x").
		Set("y = EXCLUDED.y").
		Set("z = EXCLUDED.z").
		Exec(execCtx)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("batch update timed out", "count", len(batch))
		} else {
			slog.Error("failed to batch update cache", "err", err, "count", len(batch))
		}
	}

	// Return a fresh slice with the same capacity
	return make([]*assets.Cache, 0, cap(batch))
}

// worker processes tasks from the task channel.
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

// SubmitTask adds a task to the task channel.
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

// submitDBTask submits a database task to the dedicated DB worker.
func (p *Processor) submitDBTask(task DBTask) (interface{}, error) {
	responseChan := make(chan interface{}, 1)
	errorChan := make(chan error, 1)

	task.ResponseChan = responseChan
	task.ErrorChan = errorChan

	// Submit the task
	select {
	case p.dbTaskCh <- task:
		// Wait for response or error
		select {
		case response := <-responseChan:
			return response, nil
		case err := <-errorChan:
			return nil, err
		case <-time.After(10 * time.Second):
			return nil, eris.Errorf("timeout waiting for DB task response: %s", task.Type)
		}
	case <-time.After(5 * time.Second):
		return nil, eris.Errorf("timeout submitting DB task: %s", task.Type)
	}
}

// ScanFS scans the filesystem and submits tasks.
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

// WaitForCompletion waits for all tasks to complete.
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
		p.stats.LogStatus()
		slog.Warn(
			"timeout waiting for completion",
			"timeout",
			timeout.String(),
			"taskQueueSize",
			len(p.taskCh),
			"dbQueueSize",
			len(p.dbTaskCh),
			"pendingTasks",
			p.stats.TotalSubmitted.Load()-p.stats.TotalCompleted.Load(),
			"pendingDBTasks",
			p.stats.DBTasksSubmitted.Load()-p.stats.DBTasksCompleted.Load(),
		)

		return false
	}
}

// WaitForAllTasks waits for the wait group to finish
// This is different from WaitForCompletion as it waits for the actual
// wait group to be empty, not just for the done signal.
func (p *Processor) WaitForAllTasks() {
	p.wg.Wait()
}

// Close properly cleans up all resources.
func (p *Processor) Close() {
	// Signal DB worker to shutdown and flush any pending batches
	close(p.dbShutdownCh)

	// Close all channels
	close(p.taskCh)
	close(p.dbTaskCh)
	close(p.errCh)
	close(p.doneCh)

	slog.Debug("processor resources cleaned up")
}

// Errors returns the error channel.
func (p *Processor) Errors() <-chan error {
	return p.errCh
}

// Done returns the done channel.
func (p *Processor) Done() <-chan struct{} {
	return p.doneCh
}

// processTask handles processing of a single task.
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

// processAsset processes an asset file.
func (p *Processor) processAsset(ctx context.Context, task Task) error {
	var err error

	// Read file if content wasn't provided
	if task.Content == nil {
		task.Content, err = afero.ReadFile(p.fs, task.Path)
		if err != nil {
			return eris.Wrapf(err, "failed to read asset: %s", task.Path)
		}
	}

	// Calculate hash
	hash := assets.Hash(task.Content)

	// Check if asset is already cached and unchanged
	isCached, err := p.checkCache(ctx, task.Path, hash)
	if err != nil {
		return err
	}
	if isCached {
		slog.Debug("asset cached", "path", task.Path)

		return nil // Skip if already cached and unchanged
	}

	// Upload to S3
	err = p.uploadToS3(ctx, task.Path, task.Content)
	if err != nil {
		return err
	}

	// Update cache using the DB worker
	_, err = p.submitDBTask(DBTask{
		Type: DBUpdateCache,
		Data: &assets.Cache{
			Path: task.Path,
			Hash: hash,
		},
	})

	return err
}

// processDocument processes a document file.
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
	isCached, err := p.checkCache(ctx, task.Path, hash)
	if err != nil {
		return err
	}
	if isCached {
		return nil // Skip if already cached and unchanged
	}

	// Parse document
	doc := &assets.Doc{
		Path: task.Path,
		Hash: hash,
	}

	// Set default values
	err = assets.Defaults(doc)
	if err != nil {
		return eris.Wrapf(err, "failed to set defaults for document: %s", task.Path)
	}

	// Parse markdown and frontmatter
	err = p.parseMarkdown(content, doc)
	if err != nil {
		return err
	}

	// Set slug based on path
	doc.Slug = assets.Slugify(task.Path)

	// Validate document
	err = assets.Validate(task.Path, doc)
	if err != nil {
		return eris.Wrapf(err, "document validation failed: %s", task.Path)
	}

	// Save document to database through the DB worker
	err = p.saveDocument(ctx, doc)
	if err != nil {
		return err
	}

	// Update cache through the DB worker
	_, err = p.submitDBTask(DBTask{
		Type: DBUpdateCache,
		Data: &assets.Cache{
			Path: task.Path,
			Hash: hash,
		},
	})

	return err
}

// checkCache checks if a file is already cached and unchanged.
func (p *Processor) checkCache(ctx context.Context, path, hash string) (bool, error) {
	// Check in memory cache first
	if cachedHash, ok := p.memCache.Get(path); ok {
		if cachedHash == hash {
			return true, nil
		}
		slog.Debug("file has changed", "path", path)

		return false, nil
	}

	// Not in memory cache, load from database if needed
	_, err := p.submitDBTask(DBTask{
		Type: DBLoadCache,
	})
	if err != nil {
		return false, err
	}

	// Check again after loading cache
	if cachedHash, ok := p.memCache.Get(path); ok {
		if cachedHash == hash {
			return true, nil
		}
	}

	// Not in cache
	return false, nil
}

// uploadToS3 uploads a file to S3.
func (p *Processor) uploadToS3(ctx context.Context, path string, data []byte) error {
	// Use custom content type function instead of mime.TypeByExtension
	contentType := getContentType(path)

	// Use a timeout context to prevent hanging on S3 operations
	uploadCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	slog.Info("Uploading to S3", "path", path)
	_, err := p.s3Client.PutObject(uploadCtx, &s3.PutObjectInput{
		Bucket:      &p.s3Bucket,
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

// parseMarkdown parses markdown content and extracts frontmatter.
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

// saveDocument saves a document to the database.
func (p *Processor) saveDocument(ctx context.Context, doc *assets.Doc) error {
	switch {
	case strings.HasPrefix(doc.Path, assets.PostsLoc) && strings.HasSuffix(doc.Path, ".md"):
		// Convert Doc to Post
		post := &assets.Post{}
		copygen.ToPost(post, doc)

		// Submit to DB worker
		response, err := p.submitDBTask(DBTask{
			Type: DBSavePost,
			Data: post,
		})
		if err != nil {
			return err
		}

		// Schedule relationship update task with a delay to avoid database contention
		time.AfterFunc(100*time.Millisecond, func() {
			if savedPost, ok := response.(*assets.Post); ok {
				p.SubmitTask(Task{
					Type: TypeRelationship,
					Path: doc.Path,
					RelationshipFn: func(ctx context.Context) error {
						_, err := p.submitDBTask(DBTask{
							Type: DBUpdateRelationship,
							Data: struct {
								Post *assets.Post
							}{Post: savedPost},
						})

						return err
					},
				})
			}
		})

		return nil

	case strings.HasPrefix(doc.Path, assets.ProjectsLoc) && strings.HasSuffix(doc.Path, ".md"):
		// Convert Doc to Project
		project := &assets.Project{}
		copygen.ToProject(project, doc)

		// Submit to DB worker
		response, err := p.submitDBTask(DBTask{
			Type: DBSaveProject,
			Data: project,
		})
		if err != nil {
			return err
		}

		// Schedule relationship update task with a delay to avoid database contention
		time.AfterFunc(100*time.Millisecond, func() {
			if savedProject, ok := response.(*assets.Project); ok {
				p.SubmitTask(Task{
					Type: TypeRelationship,
					Path: doc.Path,
					RelationshipFn: func(ctx context.Context) error {
						_, err := p.submitDBTask(DBTask{
							Type: DBUpdateRelationship,
							Data: struct {
								Project *assets.Project
							}{Project: savedProject},
						})

						return err
					},
				})
			}
		})

		return nil

	case strings.HasPrefix(doc.Path, assets.TagsLoc) && strings.HasSuffix(doc.Path, ".md"):
		// Convert Doc to Tag
		tag := &assets.Tag{}
		copygen.ToTag(tag, doc)

		// Submit to DB worker
		response, err := p.submitDBTask(DBTask{
			Type: DBSaveTag,
			Data: tag,
		})
		if err != nil {
			return err
		}

		// Schedule relationship update task with a delay to avoid database contention
		time.AfterFunc(100*time.Millisecond, func() {
			if savedTag, ok := response.(*assets.Tag); ok {
				p.SubmitTask(Task{
					Type: TypeRelationship,
					Path: doc.Path,
					RelationshipFn: func(ctx context.Context) error {
						_, err := p.submitDBTask(DBTask{
							Type: DBUpdateRelationship,
							Data: struct {
								Tag *assets.Tag
							}{Tag: savedTag},
						})

						return err
					},
				})
			}
		})

		return nil

	default:
		return eris.Errorf("unknown document type: %s", doc.Path)
	}
}
