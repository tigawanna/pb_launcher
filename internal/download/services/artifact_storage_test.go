package services

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCleanEmptyDirs(t *testing.T) {
	baseDir := t.TempDir()
	binaryStorage := &ArtifactStorage{baseDir: baseDir}

	emptyDir := filepath.Join(baseDir, "empty-dir")
	nonEmptyDir := filepath.Join(baseDir, "non-empty-dir")
	require.NoError(t, os.Mkdir(emptyDir, 0755))
	require.NoError(t, os.Mkdir(nonEmptyDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(nonEmptyDir, "file.txt"), []byte("data"), 0644))

	require.NoError(t, binaryStorage.cleanEmptyDirs())

	_, err := os.Stat(emptyDir)
	require.True(t, os.IsNotExist(err), "empty directory should be removed")

	_, err = os.Stat(nonEmptyDir)
	require.NoError(t, err, "non-empty directory should not be removed")
}
