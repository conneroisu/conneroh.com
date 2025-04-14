// Package cache contains the cache for the development server.
package cache

import (
	"encoding/json"
	"log/slog"
	"os"
	"sync"

	"github.com/conneroisu/conneroh.com/internal/data/gen"
	"github.com/rotisserie/eris"
	"github.com/spf13/afero"
)

type MinimalEmbedding struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// Cache is the storage of previous hashes.
type Cache struct {
	HashFile string `json:"-"`

	mu     sync.Mutex        `json:"-"`
	Hashes map[string]string `json:"hashes"`

	muDocs  sync.Mutex                  `json:"-"`
	DocsEmb map[string]MinimalEmbedding `json:"docs"`
}

// Close writes the config to disk and closes the file.
func (c *Cache) Close() (err error) {
	var tmpCache struct {
		Hashes  map[string]string           `json:"hashes"`
		DocsEmb map[string]MinimalEmbedding `json:"docs"`
	}
	tmpCache.Hashes = c.Hashes
	tmpCache.DocsEmb = c.DocsEmb
	c.Hashes = nil
	body, err := json.Marshal(tmpCache)
	if err != nil {
		return
	}
	err = os.WriteFile(c.HashFile, body, 0644)
	if err != nil {
		return
	}
	return
}

// LoadCache attempts to load a cache from the specified file path.
func LoadCache(
	path string,
) (*Cache, error) {
	slog.Debug("loading cache", "path", path)
	defer slog.Debug("cache loaded", "path", path)
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Create a new cache if the file doesn't exist
			return &Cache{
				HashFile: path,
				Hashes:   make(map[string]string),
				DocsEmb:  make(map[string]MinimalEmbedding),
			}, nil
		}
		return nil, err
	}
	defer f.Close()
	data, err := afero.ReadAll(f)
	if err != nil {
		return nil, eris.Wrapf(
			err,
			"failed to read cache from %s",
			path,
		)
	}
	var cache Cache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, eris.Wrapf(
			err,
			"failed to unmarshal cache from %s",
			path,
		)
	}
	cache.HashFile = path
	if cache.Hashes == nil {
		cache.Hashes = make(map[string]string)
	}
	if cache.DocsEmb == nil {
		cache.DocsEmb = make(map[string]MinimalEmbedding)
	}

	return &cache, nil
}

// Marshal returns the JSON representation of the cache object.
func (c *Cache) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

// Unmarshal unmarshals the JSON representation of the cache object.
func (c *Cache) Unmarshal(data []byte) error {
	return json.Unmarshal(data, c)
}

// Set sets the hash for the provided path.
func (c *Cache) Set(path string, hash string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Hashes[path] = hash
}

// SetPath sets the hash for the provided path.
func (c *Cache) SetPath(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	hash := Hash([]byte(path))
	c.Hashes[path] = hash
}

// Get returns the hash for the provided path.
func (c *Cache) Get(path string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	hash, ok := c.Hashes[path]
	return hash, ok
}

// SetDoc sets the vec data for the provided path.
func (c *Cache) SetDoc(path string, emb gen.Embedded) {
	c.muDocs.Lock()
	defer c.muDocs.Unlock()
	c.DocsEmb[path] = struct {
		X float64 "json:\"x\""
		Y float64 "json:\"y\""
		Z float64 "json:\"z\""
	}{
		X: emb.X,
		Y: emb.Y,
		Z: emb.Z,
	}
}

// GetDoc returns the vec data for the provided path.
func (c *Cache) GetDoc(path string) (MinimalEmbedding, bool) {
	c.muDocs.Lock()
	defer c.muDocs.Unlock()
	emb, ok := c.DocsEmb[path]
	return emb, ok
}
