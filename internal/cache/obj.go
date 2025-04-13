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

// Cache is the storage of previous hashes.
type Cache struct {
	HashFile string `json:"-"`

	mu     sync.Mutex        `json:"-"`
	Hashes map[string]string `json:"hashes"`

	muDocs    sync.Mutex              `json:"-"`
	DocHashes map[string]string       `json:"docHashes"`
	Docs      map[string]gen.Embedded `json:"docs"`
}

// Close writes the config to disk and closes the file.
func (c *Cache) Close() (err error) {
	var tmpCache struct {
		Hashes map[string]string `json:"hashes"`
	}
	tmpCache.Hashes = c.Hashes
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
				HashFile:  path,
				Hashes:    make(map[string]string),
				Docs:      make(map[string]gen.Embedded),
				DocHashes: make(map[string]string),
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

// Get returns the hash for the provided path.
func (c *Cache) Get(path string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	hash, ok := c.Hashes[path]
	return hash, ok
}

// SetDoc sets the hash for the provided path.
func (c *Cache) SetDoc(path string, emb gen.Embedded) {
	c.muDocs.Lock()
	defer c.muDocs.Unlock()
	c.DocHashes[path] = Hash([]byte(emb.RawContent))
	c.Docs[path] = emb
}

// OldDoc checks if the hash for the provided path is the same as the
// provided hash.
func (c *Cache) OldDoc(path string, content []byte) bool {
	c.muDocs.Lock()
	defer c.muDocs.Unlock()
	hash := Hash(content)
	known, ok := c.DocHashes[path]
	if ok && known == hash {
		return true
	}
	return false
}

// GetDoc returns the hash for the provided path.
func (c *Cache) GetDoc(path string) (gen.Embedded, bool) {
	c.muDocs.Lock()
	defer c.muDocs.Unlock()
	emb, ok := c.Docs[path]
	return emb, ok
}
