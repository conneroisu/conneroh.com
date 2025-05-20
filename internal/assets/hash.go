package assets

import (
	"context"
	//nolint:gosec
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"log/slog"
	"os"
	"sort"
	"strings"

	"github.com/rotisserie/eris"
	"github.com/spf13/afero"
	"github.com/uptrace/bun"
)

// ComputeHash generates an MD5 hash of the given content.
// Note: MD5 is not cryptographically secure; only use for content fingerprinting.
func ComputeHash(content []byte) string {
	//nolint:gosec
	sum := md5.Sum(content)

	return hex.EncodeToString(sum[:])
}

// MatchItem takes a path and returns a DirMatchItem.
func MatchItem(fs afero.Fs, path string) (DirMatchItem, error) {
	slog.Debug("assets.MatchItem()", "path", path)
	defer slog.Debug("assets.MatchItem() done", "path", path)
	content, err := afero.ReadFile(fs, path)
	if err != nil {
		return DirMatchItem{}, err
	}

	return DirMatchItem{
		Path:    path,
		Content: string(content),
		Hash:    ComputeHash(content),
	}, nil
}

// DirMatchItem contains the path and content of a file.
type DirMatchItem struct {
	Path    string
	Content string
	Hash    string
}

// HashDirMatch takes an fs, path, and a db.
//
// It returns a slice of paths if the hash of the directory does not match the
// hash in the database.
func HashDirMatch(
	ctx context.Context,
	fs afero.Fs,
	path string,
	db *bun.DB,
) ([]DirMatchItem, error) {
	slog.Debug("assets.HashDirMatch()", "path", path)
	var files []string
	var pathMap []DirMatchItem
	err := afero.Walk(
		fs,
		path,
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, p)
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	// Sort files for consistent ordering
	sort.Strings(files)

	// Create a concatenated string of paths and hashes, then hash that
	var hashInputBuilder strings.Builder

	for _, file := range files {
		var matchItem DirMatchItem
		matchItem, err = MatchItem(fs, file)
		if err != nil {
			return nil, err
		}
		pathMap = append(pathMap, matchItem)

		_, err = hashInputBuilder.WriteString(matchItem.Hash)
		if err != nil {
			return nil, err
		}
	}

	dirHash := ComputeHash([]byte(hashInputBuilder.String()))

	var dirCache Cache
	err = db.NewSelect().
		Model(&dirCache).
		Where("path = ?", path).
		Scan(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, eris.Wrap(err, "failed to check cache")
	}

	if dirCache.Hash == dirHash {
		return nil, nil
	}

	var changedFiles []DirMatchItem
	for _, value := range pathMap {
		err = db.NewSelect().
			Model(&dirCache).
			Where("path = ?", value.Path).
			Scan(ctx)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, eris.Wrap(err, "failed to check cache")
		}
		if dirCache.Hash != value.Hash {
			changedFiles = append(changedFiles, value)
		}
	}

	for _, item := range changedFiles {
		_, err = db.NewInsert().
			Model(&Cache{
				Path: item.Path,
				Hash: item.Hash,
			}).
			On("CONFLICT (path) DO UPDATE").
			Set("hashed = EXCLUDED.hashed").
			Exec(ctx)
		if err != nil {
			return nil, eris.Wrap(err, "failed to update cache")
		}
	}

	return changedFiles, nil
}
