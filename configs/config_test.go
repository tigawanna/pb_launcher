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
	originalApiDomain := os.Getenv("API_DOMAIN")

	os.Unsetenv("ADDRESS")
	os.Unsetenv("RELEASE_SYNC_INTERVAL")
	os.Unsetenv("DOWNLOAD_DIR")
	os.Setenv("API_DOMAIN", "localhost")

	config, _ := configs.ReadConfigs()
	assert.Equal(t, "0.0.0.0:7090", config.HttpAddr)
	assert.Equal(t, 10*time.Minute, config.ReleaseSyncInterval)
	assert.Equal(t, "./downloads", config.DownloadDir)

	os.Setenv("ADDRESS", "127.0.0.1:8080")
	os.Setenv("RELEASE_SYNC_INTERVAL", "15m")
	os.Setenv("DOWNLOAD_DIR", "/data/releases")
	os.Setenv("API_DOMAIN", "localhost")

	config, _ = configs.ReadConfigs()
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

	if originalApiDomain != "" {
		os.Setenv("API_DOMAIN", originalApiDomain)
	} else {
		os.Unsetenv("API_DOMAIN")
	}
}
