// Package cache contains the cache for the development server.
package cache

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/rotisserie/eris"
	"github.com/spf13/afero"
)

// Cache is the storage of previous hashes.
type Cache struct {
	fs       afero.Fs `json:"-"`
	HashFile string   `json:"-"`
	Hashes   map[string]string
}

// Close writes the config to disk and closes the file.
func (c *Cache) Close() (err error) {
	body, err := json.Marshal(c)
	if err != nil {
		return
	}
	err = afero.WriteFile(c.fs, c.HashFile, body, 0644)
	if err != nil {
		return
	}
	return
}

// LoadCache attempts to load a cache from the specified file path.
func LoadCache(
	fs afero.Fs,
	path string,
) (*Cache, error) {
	slog.Debug("loading cache", "path", path)
	defer slog.Debug("cache loaded", "path", path)
	f, err := fs.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			// Create a new cache if the file doesn't exist
			return &Cache{
				HashFile: path,
				Hashes:   make(map[string]string),
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
	cache.fs = fs

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
