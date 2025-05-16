package configs_test

import (
	"os"
	"pb_luncher/configs"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadConfigs(t *testing.T) {
	originalAddress := os.Getenv("ADDRESS")
	originalSyncInterval := os.Getenv("RELEASE_SYNC_INTERVAL")
	originalRepo := os.Getenv("GITHUB_REPOSITORY")
	originalPattern := os.Getenv("RELEASE_FILE_PATTERN")
	originalDownloadDir := os.Getenv("DOWNLOAD_DIR")

	os.Unsetenv("ADDRESS")
	os.Unsetenv("RELEASE_SYNC_INTERVAL")
	os.Unsetenv("GITHUB_REPOSITORY")
	os.Unsetenv("RELEASE_FILE_PATTERN")
	os.Unsetenv("DOWNLOAD_DIR")

	config := configs.ReadConfigs()
	assert.Equal(t, "0.0.0.0:7090", config.HttpAddr)
	assert.Equal(t, 10*time.Minute, config.ReleaseSyncInterval)
	assert.Equal(t, "pocketbase/pocketbase", config.GithubRepository)
	assert.Equal(t, "./downloads", config.DownloadDir)
	assert.True(t, config.ReleaseFilePattern.MatchString("pocketbase_0.28.1_linux_amd64.zip"))

	os.Setenv("ADDRESS", "127.0.0.1:8080")
	os.Setenv("RELEASE_SYNC_INTERVAL", "15m")
	os.Setenv("GITHUB_REPOSITORY", "custom/repo")
	os.Setenv("RELEASE_FILE_PATTERN", `myapp_.+_linux_amd64\.zip`)
	os.Setenv("DOWNLOAD_DIR", "/data/releases")

	config = configs.ReadConfigs()
	assert.Equal(t, "127.0.0.1:8080", config.HttpAddr)
	assert.Equal(t, 15*time.Minute, config.ReleaseSyncInterval)
	assert.Equal(t, "custom/repo", config.GithubRepository)
	assert.Equal(t, "/data/releases", config.DownloadDir)
	assert.True(t, config.ReleaseFilePattern.MatchString("myapp_1.0.0_linux_amd64.zip"))

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

	if originalRepo != "" {
		os.Setenv("GITHUB_REPOSITORY", originalRepo)
	} else {
		os.Unsetenv("GITHUB_REPOSITORY")
	}

	if originalPattern != "" {
		os.Setenv("RELEASE_FILE_PATTERN", originalPattern)
	} else {
		os.Unsetenv("RELEASE_FILE_PATTERN")
	}

	if originalDownloadDir != "" {
		os.Setenv("DOWNLOAD_DIR", originalDownloadDir)
	} else {
		os.Unsetenv("DOWNLOAD_DIR")
	}
}
