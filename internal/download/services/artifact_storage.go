package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"pb_luncher/configs"
	"pb_luncher/internal/download/domain/services"
	"runtime"

	"github.com/hashicorp/go-version"
	"github.com/wailsapp/mimetype"
)

type ArtifactStorage struct {
	baseDir string
}

var _ services.ArtifactStorage = (*ArtifactStorage)(nil)

func NewArtifactStorage(c *configs.Configs) *ArtifactStorage {
	return &ArtifactStorage{
		baseDir: c.DownloadDir,
	}
}

func (b *ArtifactStorage) cleanEmptyDirs() error {
	entries, err := os.ReadDir(b.baseDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("failed to read base directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirPath := filepath.Join(b.baseDir, entry.Name())
		subEntries, err := os.ReadDir(dirPath)
		if err != nil {
			slog.Error("failed to read subdirectory", "dir", dirPath, "error", err)
			continue
		}

		if len(subEntries) == 0 {
			if err := os.Remove(dirPath); err != nil {
				slog.Warn("failed to remove empty subdirectory", "dir", dirPath, "error", err)
			}
		}
	}

	return nil
}

// Versions implements services.ArtifactStorage.
func (b *ArtifactStorage) Versions(ctx context.Context) ([]*version.Version, error) {
	if err := b.cleanEmptyDirs(); err != nil {
		slog.Warn("failed to clean empty directories", "error", err)
	}

	entries, err := os.ReadDir(b.baseDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		slog.Error("failed to read base directory", "error", err, "baseDir", b.baseDir)
		return nil, fmt.Errorf("failed to read base directory: %w", err)
	}

	var versions []*version.Version
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		v, err := version.NewVersion(entry.Name())
		if err != nil {
			slog.Warn("invalid version directory", "dir", entry.Name(), "error", err)
			continue
		}

		versions = append(versions, v)
	}
	return versions, nil
}

// Save stores the provided binary data from the given reader at the specified relative path,
// creating necessary directories and setting executable permissions for binary files.
func (b *ArtifactStorage) Save(ctx context.Context, relativePath string, reader io.Reader) (string, error) {
	binaryPath := filepath.Join(b.baseDir, relativePath)

	if err := os.MkdirAll(filepath.Dir(binaryPath), 0755); err != nil {
		slog.Error("failed to create directory", "error", err, "path", filepath.Dir(binaryPath))
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(binaryPath)
	if err != nil {
		slog.Error("failed to create binary file", "error", err, "path", binaryPath)
		return "", fmt.Errorf("failed to create binary file: %w", err)
	}
	defer file.Close()

	// Copy data to the file
	if _, err := io.Copy(file, reader); err != nil {
		slog.Error("failed to write binary file", "error", err, "path", binaryPath)
		return "", fmt.Errorf("failed to write binary file: %w", err)
	}

	// Reopen file to detect MIME type
	file.Seek(0, 0)
	buffer := make([]byte, 1024)
	n, _ := file.Read(buffer)
	mime := mimetype.Detect(buffer[:n])

	isExecutable := false
	if runtime.GOOS == "windows" {
		isExecutable = mime.Is("application/vnd.microsoft.portable-executable")
	} else {
		isExecutable = mime.Is("application/x-executable") ||
			mime.Is("application/x-sharedlib") ||
			mime.Is("application/octet-stream")
	}

	if isExecutable {
		if err := os.Chmod(binaryPath, 0755); err != nil {
			slog.Error("failed to set execution permissions", "error", err, "path", binaryPath)
			return "", fmt.Errorf("failed to set execution permissions: %w", err)
		}
	}
	return binaryPath, nil
}
