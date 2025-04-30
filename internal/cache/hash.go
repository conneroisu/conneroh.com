package cache

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/spf13/afero"
)

// WorkCache is a map of work directories to their hashes.
type WorkCache map[string]string

// HashCache represents a persistent cache of path hashes.
type HashCache struct {
	mu                 sync.RWMutex
	Paths              WorkCache
	CachePath          string
	fs                 afero.Fs
	lastComputedHashes WorkCache // Internal cache of last computed hashes
}

// NewHashCache creates a new cache with an optional path for persistence.
func NewHashCache(fs afero.Fs, cachePath string) (*HashCache, error) {
	cache := &HashCache{
		Paths:              make(WorkCache),
		CachePath:          cachePath,
		fs:                 fs,
		lastComputedHashes: nil,
	}

	if cachePath != "" {
		err := cache.Load()
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}

	return cache, nil
}

// Load reads the cache from disk.
func (c *HashCache) Load() error {
	if c.CachePath == "" {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := afero.ReadFile(c.fs, c.CachePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &c.Paths)
}

// Save writes the cache to disk.
func (c *HashCache) Save() error {
	if c.CachePath == "" {
		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	data, err := json.Marshal(c.Paths)
	if err != nil {
		return err
	}

	return afero.WriteFile(c.fs, c.CachePath, data, 0644)
}

// Close ensures the cache is saved before exiting.
func (c *HashCache) Close() error {
	return c.Save()
}

// SplitGlobPath splits a glob path into a base directory and pattern.
// For example: "./src/*.go" -> "./src", "*.go"
func SplitGlobPath(globPath string) (string, string) {
	// Check if the path contains glob metacharacters
	if !strings.ContainsAny(globPath, "*?[{") {
		return globPath, ""
	}

	// Find the last directory separator before a glob character
	lastSep := -1
	for i := range len(globPath) {
		if globPath[i] == '/' || globPath[i] == '\\' {
			lastSep = i
		}
		if strings.ContainsRune("*?[{", rune(globPath[i])) {
			break
		}
	}

	// If we found a separator, split the path
	if lastSep >= 0 {
		dir := globPath[:lastSep]
		pattern := globPath[lastSep+1:]

		// Handle empty directory (should be "." for current directory)
		if dir == "" {
			dir = "."
		}

		return dir, pattern
	}

	// If we didn't find a separator before a glob character,
	// the pattern is in the current directory
	return ".", globPath
}

// hashFile computes the MD5 hash of a file.
func hashFile(fs afero.Fs, path string) (string, error) {
	file, err := fs.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// hashGlobPath computes a hash for a path that may contain a glob pattern.
func hashGlobPath(fs afero.Fs, globPath string) (string, error) {
	dir, pattern := SplitGlobPath(globPath)
	return hashDir(fs, dir, pattern)
}

// hashDir computes a hash for a directory by hashing its structure and content.
func hashDir(fs afero.Fs, path string, globPattern string) (string, error) {
	var files []string

	// If we have a glob pattern, use it to match files
	if globPattern != "" {
		// Use afero.Glob if available, otherwise implement our own
		matches, err := aferoGlob(fs, filepath.Join(path, globPattern))
		if err != nil {
			return "", err
		}
		files = matches
	} else {
		// Otherwise walk the entire directory
		err := afero.Walk(fs, path, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, p)
			}
			return nil
		})
		if err != nil {
			return "", err
		}
	}

	// Sort files for consistent ordering
	sort.Strings(files)

	// Hash each file and combine
	h := md5.New()
	for _, file := range files {
		relPath, err := filepath.Rel(path, file)
		if err != nil {
			return "", err
		}

		// Include the relative path in the hash
		_, err = io.WriteString(h, relPath)
		if err != nil {
			return "", err
		}

		fileHash, err := hashFile(fs, file)
		if err != nil {
			return "", err
		}

		_, err = io.WriteString(h, fileHash)
		if err != nil {
			return "", err
		}
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// aferoGlob implements a simple glob for afero filesystems
func aferoGlob(fs afero.Fs, pattern string) ([]string, error) {
	dir, filePattern := filepath.Split(pattern)
	if dir == "" {
		dir = "."
	}

	// Check if the directory exists
	dirInfo, err := fs.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !dirInfo.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dir)
	}

	// Read all files in the directory
	dirFile, err := fs.Open(dir)
	if err != nil {
		return nil, err
	}
	defer dirFile.Close()

	dirList, err := dirFile.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	// Match files against the pattern
	var matches []string
	for _, fileName := range dirList {
		matched, err := filepath.Match(filePattern, fileName)
		if err != nil {
			return nil, err
		}
		if matched {
			fullPath := filepath.Join(dir, fileName)
			// Check if it's a file (not a directory)
			fileInfo, err := fs.Stat(fullPath)
			if err != nil {
				return nil, err
			}
			if !fileInfo.IsDir() {
				matches = append(matches, fullPath)
			}
		}
	}

	return matches, nil
}

// ComputeGlobHashes calculates hashes for each glob path and returns a mapping.
func ComputeGlobHashes(fs afero.Fs, globPaths []string) (WorkCache, error) {
	cache := make(WorkCache)
	for _, globPath := range globPaths {
		hash, err := hashGlobPath(fs, globPath)
		if err != nil {
			return nil, err
		}

		cache[globPath] = hash
	}

	return cache, nil
}

// HasGlobChanges checks if any glob paths have changed compared to the cached hashes.
func (c *HashCache) HasGlobChanges(globPaths []string) (bool, error) {
	currentHashes, err := ComputeGlobHashes(c.fs, globPaths)
	if err != nil {
		return false, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	for path, hash := range currentHashes {
		cachedHash, exists := c.Paths[path]
		if !exists || cachedHash != hash {
			// Store the current hashes for later use in UpdateGlobHashes
			c.mu.RUnlock()
			c.mu.Lock()
			c.lastComputedHashes = currentHashes
			c.mu.Unlock()
			c.mu.RLock()
			return true, nil
		}
	}

	return false, nil
}

// UpdateGlobHashes stores new hashes for the given glob paths.
// If useLastComputed is true, it will use the hashes from the last HasGlobChanges call.
func (c *HashCache) UpdateGlobHashes(globPaths []string, useLastComputed bool) error {
	var newHashes WorkCache
	var err error

	if useLastComputed && c.lastComputedHashes != nil {
		c.mu.RLock()
		// Check if we have pre-computed hashes for these exact paths
		allPathsFound := true
		for _, path := range globPaths {
			if _, exists := c.lastComputedHashes[path]; !exists {
				allPathsFound = false
				break
			}
		}

		if allPathsFound {
			newHashes = c.lastComputedHashes
			c.mu.RUnlock()
		} else {
			c.mu.RUnlock()
			// Some paths were not in lastComputedHashes, recompute all
			newHashes, err = ComputeGlobHashes(c.fs, globPaths)
			if err != nil {
				return err
			}
		}
	} else {
		// No pre-computed hashes or not requested to use them
		newHashes, err = ComputeGlobHashes(c.fs, globPaths)
		if err != nil {
			return err
		}
	}

	c.mu.Lock()
	maps.Copy(c.Paths, newHashes)
	c.lastComputedHashes = nil // Clear cache after using it
	c.mu.Unlock()

	return c.Save()
}

// OnGlobChanges executes the given function if any glob paths have changed.
// If changes are detected, the function is executed and the cache is updated.
func (c *HashCache) OnGlobChanges(globPaths []string, fn func() error) error {
	changed, err := c.HasGlobChanges(globPaths)
	if err != nil {
		return err
	}

	if !changed {
		return nil
	}

	// Execute the function
	if err := fn(); err != nil {
		return err
	}

	// Update the hashes after successful execution using the already computed hashes
	return c.UpdateGlobHashes(globPaths, true)
}

// OnGlobPathChanges is a simplified version of OnGlobChanges that uses a WorkCache directly.
func OnGlobPathChanges(fs afero.Fs, workCache WorkCache, globPaths []string, fn func() error) (WorkCache, error) {
	currentHashes, err := ComputeGlobHashes(fs, globPaths)
	if err != nil {
		return workCache, err
	}

	// Check for changes
	changed := false
	for path, hash := range currentHashes {
		cachedHash, exists := workCache[path]
		if !exists || cachedHash != hash {
			changed = true
			break
		}
	}

	if !changed {
		return workCache, nil
	}

	// Execute the function
	if err := fn(); err != nil {
		return workCache, err
	}

	// Update the hashes after successful execution
	if workCache == nil {
		workCache = make(WorkCache)
	}
	maps.Copy(workCache, currentHashes)

	return workCache, nil
}
