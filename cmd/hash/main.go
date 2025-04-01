package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"github.com/conneroisu/conneroh.com/internal/cache"
)

var (
	// Command line flags
	dirPath         = flag.String("dir", "", "Path to the directory to check for changes")
	verbose         = flag.Bool("v", false, "Enable verbose output")
	excludePatterns = flag.String("exclude", "", "Comma-separated list of glob patterns to exclude")
	hashFilePath    = flag.String("cache", "", "Path to the cache file (defaults to .dir_hash.json in the directory)")
)

const defaultHashFileName = "config.json"

func main() {
	// Setup signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to listen for signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Handle signals in a separate goroutine
	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v, shutting down...", sig)
		cancel()
		// Give a short grace period for cleanup and then force exit
		go func() {
			<-sigChan
			log.Println("Received second signal, forcing exit")
			os.Exit(1)
		}()
	}()

	if err := run(ctx); err != nil {
		if err == context.Canceled {
			log.Println("Operation was canceled")
			os.Exit(1)
		}
		log.Fatalf("Error: %v", err)
	}
}

func run(ctx context.Context) error {
	flag.Parse()

	// Get directory path from flag or positional argument
	dirPathValue := *dirPath
	if dirPathValue == "" {
		if flag.NArg() > 0 {
			dirPathValue = flag.Arg(0)
		} else {
			return fmt.Errorf("directory path is required. Use -dir flag or provide as positional argument")
		}
	}

	// Process exclude patterns
	var excludes []string
	if *excludePatterns != "" {
		excludes = strings.Split(*excludePatterns, ",")
		if *verbose {
			fmt.Println("Excluding patterns:", excludes)
		}
	}

	// Check if directory exists
	info, err := os.Stat(dirPathValue)
	if err != nil {
		return fmt.Errorf("cannot access directory %s: %w", dirPathValue, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", dirPathValue)
	}

	// Check for context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Use default hash file path if not specified
	hashFilePathValue := *hashFilePath
	if hashFilePathValue == "" {
		hashFilePathValue = filepath.Join(dirPathValue, defaultHashFileName)
	}

	// Load or create cache
	cache, err := cache.LoadCache(hashFilePathValue)
	if err != nil {
		return fmt.Errorf("error loading cache: %w", err)
	}

	// Set hash file in cache if it's a new cache
	cache.HashFile = hashFilePathValue

	// Check for context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Calculate the hash of the directory
	currentHash, err := calculateDirectoryHash(ctx, dirPathValue, excludes, *verbose)
	if err != nil {
		if err == context.Canceled {
			return err
		}
		return fmt.Errorf("error calculating directory hash: %w", err)
	}

	if *verbose {
		fmt.Printf("Current hash of %s: %s\n", dirPathValue, currentHash)
	}

	// Get the previous hash for this directory
	previousHash := cache.DirHash

	if previousHash == "" {
		if *verbose {
			fmt.Println("No previous hash found")
		}
		// First run, update the cache and exit with code 0
		cache.DirHash = currentHash
		if err := cache.Close(); err != nil {
			return fmt.Errorf("error writing cache: %w", err)
		}
		fmt.Println("Initial hash created")
		os.Exit(0)
	}

	// Compare hashes
	if currentHash != previousHash {
		if *verbose {
			fmt.Printf("Changes detected in %s\n", dirPathValue)
			fmt.Printf("Previous hash: %s\n", previousHash)
			fmt.Printf("Current hash: %s\n", currentHash)
		} else {
			fmt.Printf("Changes detected in %s\n", dirPathValue)
		}

		// Update the cache
		cache.DirHash = currentHash
		if err := cache.Close(); err != nil {
			return fmt.Errorf("error writing cache: %w", err)
		}

		// Exit with code 1 to indicate changes were detected
		os.Exit(1)
	} else {
		fmt.Printf("No changes detected in %s\n", dirPathValue)
		// Exit with code 0 to indicate no changes
		os.Exit(0)
	}

	return nil
}

// calculateDirectoryHash computes a hash of all files in the directory using MD5
func calculateDirectoryHash(ctx context.Context, dirPath string, excludes []string, verbose bool) (string, error) {
	var fileInfos []string
	walkErr := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		// Check for context cancellation periodically
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err != nil {
			return err
		}
		// Skip the hash file itself
		if d.Name() == defaultHashFileName {
			return nil
		}
		// Skip directories themselves
		if d.IsDir() {
			return nil
		}
		// Check if this path should be excluded
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}
		var matched bool
		for _, pattern := range excludes {
			matched, err = filepath.Match(pattern, relPath)
			if err != nil {
				return err
			}
			if matched {
				if verbose {
					fmt.Printf("Excluding: %s\n", relPath)
				}
				return nil
			}
		}
		// Get file info
		info, err := d.Info()
		if err != nil {
			return err
		}
		// Calculate file hash
		fileHash, err := calculateFileHash(ctx, path)
		if err != nil {
			return err
		}
		// Store relative path, mode, size, modification time and hash
		fileData := fmt.Sprintf("%s %s %d %s %s",
			relPath,
			info.Mode().String(),
			info.Size(),
			info.ModTime().String(),
			fileHash)
		fileInfos = append(fileInfos, fileData)
		return nil
	})

	if walkErr != nil {
		if walkErr == context.Canceled {
			return "", walkErr
		}
		return "", walkErr
	}

	// Check context after the walk but before sorting
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	// Sort file infos to ensure consistent ordering
	sort.Strings(fileInfos)

	// Combine and hash all file information
	hasher := md5.New()
	for _, fileInfo := range fileInfos {
		if verbose {
			fmt.Println(fileInfo)
		}
		_, err := io.WriteString(hasher, fileInfo+"\n")
		if err != nil {
			return "", err
		}
	}

	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash), nil
}

// calculateFileHash computes the MD5 hash of a single file
func calculateFileHash(ctx context.Context, filePath string) (string, error) {
	// Check context before opening file
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := md5.New()

	// Create a buffer for reading
	buffer := make([]byte, 32*1024)

	for {
		// Check context before reading
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}

		n, err := file.Read(buffer)
		if n > 0 {
			hasher.Write(buffer[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
	}

	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash), nil
}
