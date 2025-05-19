package cache

import (
	"context"
	"log/slog"
	"time"

	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/rotisserie/eris"
)

// dbWorker is the dedicated goroutine for handling all database operations.
func (p *Processor) dbWorker(ctx context.Context) {
	slog.Debug("starting database worker")
	defer slog.Debug("database worker stopped")

	// Load cache initially to avoid repeated DB queries
	if err := p.memCache.loadFromDB(ctx, p.db); err != nil {
		select {
		case <-ctx.Done():
			return
		case p.errCh <- eris.Wrap(err, "db worker: failed to load initial cache"):
		}
	}

	// Maintain a periodic batch flush timer for cache updates
	cacheBatch := make([]*assets.Cache, 0, 50)
	flushTicker := time.NewTicker(500 * time.Millisecond)
	defer flushTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Flush any remaining cache updates
			if len(cacheBatch) > 0 {
				p.flushCacheBatch(ctx, cacheBatch)
			}

			return
		case <-p.dbShutdownCh:
			// Flush any remaining cache updates
			if len(cacheBatch) > 0 {
				p.flushCacheBatch(ctx, cacheBatch)
			}

			return
		case <-flushTicker.C:
			// Periodic flush of cache updates
			if len(cacheBatch) > 0 {
				newBatch := p.flushCacheBatch(ctx, cacheBatch)
				cacheBatch = newBatch
			}
		case task, ok := <-p.dbTaskCh:
			if !ok {
				return // Channel closed
			}

			p.stats.RecordDBTaskSubmitted()

			var err error
			var response interface{}

			// Process the database task
			switch task.Type {
			case DBUpdateCache:
				cache, ok := task.Data.(*assets.Cache)
				if !ok {
					err = eris.New("invalid data type for DBUpdateCache")

					break
				}

				// Add to batch
				cacheCopy := *cache // Make a copy to avoid race conditions
				cacheBatch = append(cacheBatch, &cacheCopy)

				// If batch is large enough, flush immediately
				if len(cacheBatch) >= 50 {
					cacheBatch = p.flushCacheBatch(ctx, cacheBatch)
				}

			case DBLoadCache:
				err = p.memCache.loadFromDB(ctx, p.db)

			case DBSavePost:
				post, ok := task.Data.(*assets.Post)
				if !ok {
					err = eris.New("invalid data type for DBSavePost")

					break
				}
				err = assets.UpsertPost(ctx, p.db, post)
				response = post

			case DBSaveProject:
				project, ok := task.Data.(*assets.Project)
				if !ok {
					err = eris.New("invalid data type for DBSaveProject")

					break
				}
				err = assets.UpsertProject(ctx, p.db, project)
				response = project

			case DBSaveTag:
				tag, ok := task.Data.(*assets.Tag)
				if !ok {
					err = eris.New("invalid data type for DBSaveTag")

					break
				}
				err = assets.UpsertTag(ctx, p.db, tag)
				response = tag

			case DBFindTag:
				req, ok := task.Data.(struct {
					Slug   string
					Origin string
				})
				if !ok {
					err = eris.New("invalid data type for DBFindTag")

					break
				}
				var tag *assets.Tag
				tag, err = p.dbFindTagBySlug(ctx, req.Slug, req.Origin)
				response = tag

			case DBFindPost:
				req, ok := task.Data.(struct {
					Slug   string
					Origin string
				})
				if !ok {
					err = eris.New("invalid data type for DBFindPost")

					break
				}
				var post *assets.Post
				post, err = p.dbFindPostBySlug(ctx, req.Slug, req.Origin)
				response = post

			case DBFindProject:
				req, ok := task.Data.(struct {
					Slug   string
					Origin string
				})
				if !ok {
					err = eris.New("invalid data type for DBFindProject")

					break
				}
				var project *assets.Project
				project, err = p.dbFindProjectBySlug(ctx, req.Slug, req.Origin)
				response = project

			case DBUpdateRelationship:
				switch relationData := task.Data.(type) {
				case struct {
					Post *assets.Post
				}:
					err = p.UpsertPostRelationShips(ctx, relationData.Post)
				case struct {
					Project *assets.Project
				}:
					err = p.dbUpdateProjectRelationships(ctx, relationData.Project)
				case struct {
					Tag *assets.Tag
				}:
					err = p.dbUpdateTagRelationships(ctx, relationData.Tag)
				default:
					err = eris.New("invalid data type for DBUpdateRelationship")
				}
			default:
				err = eris.Errorf("unknown database task type: %s", task.Type)
			}

			// Send response or error
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case task.ErrorChan <- err:
				}
			} else if task.ResponseChan != nil {
				select {
				case <-ctx.Done():
					return
				case task.ResponseChan <- response:
				}
			}

			p.stats.RecordDBTaskCompleted()
		}
	}
}
