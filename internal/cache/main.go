// Package cache contains the cache for the development server.
package cache

import (
	"encoding/json"
	"os"
)

// Cache is the storage of previous hashes.
type Cache struct {
	HashFile string `json:"-"`
	Hashes   map[string]string
}

// Close writes the config to disk and closes the file.
func (c *Cache) Close() (err error) {
	body, err := json.Marshal(c)
	if err != nil {
		return
	}
	err = os.WriteFile(c.HashFile, body, 0644)
	if err != nil {
		return
	}
	return
}

// LoadCache attempts to load a cache from the specified file path
func LoadCache(hashFilePath string) (*Cache, error) {
	data, err := os.ReadFile(hashFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create a new cache if the file doesn't exist
			return &Cache{
				HashFile: hashFilePath,
				Hashes:   make(map[string]string),
			}, nil
		}
		return nil, err
	}

	var cache Cache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}
	cache.HashFile = hashFilePath

	return &cache, nil
}
