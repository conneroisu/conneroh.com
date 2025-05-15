package cache_test

import (
	"testing"

	"github.com/conneroisu/conneroh.com/internal/cache"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashCacheWithTestdata(t *testing.T) {
	// Create an in-memory filesystem
	memFs := afero.NewMemMapFs()

	// Create test directory structure
	require.NoError(t, memFs.MkdirAll("testdata/src", 0755))
	require.NoError(t, memFs.MkdirAll("testdata/static", 0755))

	// Create test files
	goContent := []byte(`package src

func Hello() {
	println("Hello, World!")
}`)
	htmlContent := []byte(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
</head>
<body>
    <h1>Hello, World!</h1>
</body>
</html>`)

	require.NoError(t, afero.WriteFile(memFs, "testdata/src/hello.go", goContent, 0644))
	require.NoError(t, afero.WriteFile(memFs, "testdata/static/index.html", htmlContent, 0644))

	// Create a temporary cache file path
	cacheFile := "test-cache.json"

	// Create a new HashCache using the in-memory filesystem
	hashCache, err := cache.NewHashCache(memFs, cacheFile)
	require.NoError(t, err)
	defer hashCache.Close()

	// Define paths to monitor
	globPaths := []string{
		"testdata/src/*.go",
		"testdata/static/*.html",
	}

	// Test 1: Initial run should detect changes (empty cache)
	changeDetected, err := hashCache.HasGlobChanges(globPaths)
	require.NoError(t, err)
	assert.True(t, changeDetected, "Should detect changes on first run with empty cache")

	// Update the cache with current hashes
	err = hashCache.UpdateGlobHashes(globPaths, true)
	require.NoError(t, err)

	// Test 2: Second run should not detect changes
	changeDetected, err = hashCache.HasGlobChanges(globPaths)
	require.NoError(t, err)
	assert.False(t, changeDetected, "Should not detect changes after updating cache")

	// Test 3: OnGlobChanges should not run function when no changes
	functionRun := false
	err = hashCache.OnGlobChanges(globPaths, func() error {
		functionRun = true

		return nil
	})
	require.NoError(t, err)
	assert.False(t, functionRun, "Function should not run when no changes detected")

	// Test 4: Modify a file and check again
	modifiedGoContent := []byte(`package src

func Hello() {
	println("Hello, Modified World!")
}`)
	require.NoError(t, afero.WriteFile(memFs, "testdata/src/hello.go", modifiedGoContent, 0644))

	// Test 5: Should detect changes after file modification
	changeDetected, err = hashCache.HasGlobChanges(globPaths)
	require.NoError(t, err)
	assert.True(t, changeDetected, "Should detect changes after file modification")

	// Test 6: OnGlobChanges should run function when changes are detected
	functionRun = false
	err = hashCache.OnGlobChanges(globPaths, func() error {
		functionRun = true

		return nil
	})
	require.NoError(t, err)
	assert.True(t, functionRun, "Function should run when changes detected")

	// Test 7: Simple WorkCache usage
	workCache := cache.WorkCache{}
	functionRun = false
	workCache, err = cache.OnGlobPathChanges(memFs, workCache, globPaths, func() error {
		functionRun = true

		return nil
	})
	require.NoError(t, err)
	assert.True(t, functionRun, "Function should run with WorkCache when changes detected")

	// Update the file again to test the last version
	finalGoContent := []byte(`package src

func Hello() {
	println("Hello, World Version 3!")
}`)
	require.NoError(t, afero.WriteFile(memFs, "testdata/src/hello.go", finalGoContent, 0644))

	// Test 8: WorkCache should detect changes again
	functionRun = false
	_, err = cache.OnGlobPathChanges(memFs, workCache, globPaths, func() error {
		functionRun = true

		return nil
	})
	require.NoError(t, err)
	assert.True(t, functionRun, "Function should run after another file change")
}

func TestHashCacheDirectoryMonitoring(t *testing.T) {
	// Create an in-memory filesystem
	memFs := afero.NewMemMapFs()

	// Create a test directory with files
	require.NoError(t, memFs.MkdirAll("testfiles", 0755))

	// Create test files
	require.NoError(t, afero.WriteFile(memFs, "testfiles/file1.txt", []byte("test content 1"), 0644))
	require.NoError(t, afero.WriteFile(memFs, "testfiles/file2.go", []byte("package main\n\nfunc main() {}\n"), 0644))

	// Create a new HashCache using the in-memory filesystem
	cacheFile := "dir-cache.json"
	hashCache, err := cache.NewHashCache(memFs, cacheFile)
	require.NoError(t, err)
	defer hashCache.Close()

	// Monitor the entire directory
	dirPath := "testfiles"

	// Test 1: Initial run should detect changes (empty cache)
	changeDetected, err := hashCache.HasGlobChanges([]string{dirPath})
	require.NoError(t, err)
	assert.True(t, changeDetected, "Should detect changes on first run with empty cache")

	// Update the cache
	err = hashCache.UpdateGlobHashes([]string{dirPath}, true)
	require.NoError(t, err)

	// Test 2: No changes after update
	changeDetected, err = hashCache.HasGlobChanges([]string{dirPath})
	require.NoError(t, err)
	assert.False(t, changeDetected, "Should not detect changes after updating cache")

	// Test 3: Add a new file
	require.NoError(t, afero.WriteFile(memFs, "testfiles/file3.txt", []byte("new file content"), 0644))

	changeDetected, err = hashCache.HasGlobChanges([]string{dirPath})
	require.NoError(t, err)
	assert.True(t, changeDetected, "Should detect changes after adding new file")

	// Test 4: Remove a file
	require.NoError(t, memFs.Remove("testfiles/file1.txt"))

	// Update cache first to recognize the new file3
	err = hashCache.UpdateGlobHashes([]string{dirPath}, true)
	require.NoError(t, err)

	changeDetected, err = hashCache.HasGlobChanges([]string{dirPath})
	require.NoError(t, err)
	assert.True(t, changeDetected, "Should detect changes after removing a file")
}
