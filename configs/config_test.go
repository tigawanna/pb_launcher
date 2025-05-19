package configs_test

import (
	"os"
	"pb_launcher/configs"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadConfigs(t *testing.T) {
	originalAddress := os.Getenv("ADDRESS")
	originalSyncInterval := os.Getenv("RELEASE_SYNC_INTERVAL")
	originalDownloadDir := os.Getenv("DOWNLOAD_DIR")

	os.Unsetenv("ADDRESS")
	os.Unsetenv("RELEASE_SYNC_INTERVAL")
	os.Unsetenv("DOWNLOAD_DIR")

	config := configs.ReadConfigs()
	assert.Equal(t, "0.0.0.0:7090", config.HttpAddr)
	assert.Equal(t, 10*time.Minute, config.ReleaseSyncInterval)
	assert.Equal(t, "./downloads", config.DownloadDir)

	os.Setenv("ADDRESS", "127.0.0.1:8080")
	os.Setenv("RELEASE_SYNC_INTERVAL", "15m")
	os.Setenv("GITHUB_REPOSITORY", "custom/repo")
	os.Setenv("RELEASE_FILE_PATTERN", `myapp_.+_linux_amd64\.zip`)
	os.Setenv("DOWNLOAD_DIR", "/data/releases")

	config = configs.ReadConfigs()
	assert.Equal(t, "127.0.0.1:8080", config.HttpAddr)
	assert.Equal(t, 15*time.Minute, config.ReleaseSyncInterval)
	assert.Equal(t, "/data/releases", config.DownloadDir)

	if originalAddress != "" {
		os.Setenv("ADDRESS", originalAddress)
	} else {
		os.Unsetenv("ADDRESS")
	}

	if originalSyncInterval != "" {
		os.Setenv("RELEASE_SYNC_INTERVAL", originalSyncInterval)
	} else {
		os.Unsetenv("RELEASE_SYNC_INTERVAL")
	}

	if originalDownloadDir != "" {
		os.Setenv("DOWNLOAD_DIR", originalDownloadDir)
	} else {
		os.Unsetenv("DOWNLOAD_DIR")
	}
}
