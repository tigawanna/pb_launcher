package services

import (
	"context"
	"fmt"
	"os"
	"path"
	"pb_launcher/configs"
	"pb_launcher/internal/launcher/domain/services"
	"regexp"
)

type BinaryFinder struct {
	targetDir string
}

var _ services.BinaryFinder = (*BinaryFinder)(nil)

func NewBinaryFinder(c *configs.Configs) *BinaryFinder {
	return &BinaryFinder{
		targetDir: c.DownloadDir,
	}
}

func (bf *BinaryFinder) FindBinary(ctx context.Context, repositoryID, version string, binaryPattern *regexp.Regexp) (string, error) {
	repoDir := path.Join(bf.targetDir, repositoryID, version)

	files, err := os.ReadDir(repoDir)
	if err != nil {
		return "", fmt.Errorf("failed to read directory %s: %w", repoDir, err)
	}

	for _, entry := range files {
		if entry.IsDir() {
			continue
		}
		fileName := entry.Name()
		if binaryPattern.MatchString(fileName) {
			return path.Join(repoDir, fileName), nil
		}
	}

	return "", fmt.Errorf("no binary found matching pattern %s in repository %s", binaryPattern.String(), repositoryID)
}
