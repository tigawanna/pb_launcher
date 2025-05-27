package configs_test

import (
	"os"
	"pb_launcher/configs"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadConfigs(t *testing.T) {
	originalBindAddress := os.Getenv("BIND_ADDRESS")
	originalBindPort := os.Getenv("BIND_PORT")
	originalSyncInterval := os.Getenv("RELEASE_SYNC_INTERVAL")
	originalDownloadDir := os.Getenv("DOWNLOAD_DIR")
	originalApiDomain := os.Getenv("PUBLIC_API_DOMAIN")
	originalDataDir := os.Getenv("DATA_DIR")

	os.Unsetenv("BIND_ADDRESS")
	os.Unsetenv("BIND_PORT")
	os.Unsetenv("RELEASE_SYNC_INTERVAL")
	os.Unsetenv("DOWNLOAD_DIR")
	os.Unsetenv("DATA_DIR")
	os.Setenv("PUBLIC_API_DOMAIN", "localhost")

	config, err := configs.ReadConfigs()
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1", config.BindAddress)
	assert.Equal(t, "8072", config.BindPort)
	assert.Equal(t, 10*time.Minute, config.ReleaseSyncInterval)
	assert.Equal(t, "./downloads", config.DownloadDir)
	assert.Equal(t, "./data", config.DataDir)
	assert.Equal(t, "localhost", config.PublicApiDomain)

	os.Setenv("BIND_ADDRESS", "0.0.0.0")
	os.Setenv("BIND_PORT", "7090")
	os.Setenv("RELEASE_SYNC_INTERVAL", "15m")
	os.Setenv("DOWNLOAD_DIR", "/data/releases")
	os.Setenv("DATA_DIR", "/data/app")
	os.Setenv("PUBLIC_API_DOMAIN", "example.test")

	config, err = configs.ReadConfigs()
	assert.NoError(t, err)
	assert.Equal(t, "0.0.0.0", config.BindAddress)
	assert.Equal(t, "7090", config.BindPort)
	assert.Equal(t, 15*time.Minute, config.ReleaseSyncInterval)
	assert.Equal(t, "/data/releases", config.DownloadDir)
	assert.Equal(t, "/data/app", config.DataDir)
	assert.Equal(t, "example.test", config.PublicApiDomain)

	if originalBindAddress != "" {
		os.Setenv("BIND_ADDRESS", originalBindAddress)
	} else {
		os.Unsetenv("BIND_ADDRESS")
	}
	if originalBindPort != "" {
		os.Setenv("BIND_PORT", originalBindPort)
	} else {
		os.Unsetenv("BIND_PORT")
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
	if originalDataDir != "" {
		os.Setenv("DATA_DIR", originalDataDir)
	} else {
		os.Unsetenv("DATA_DIR")
	}
	if originalApiDomain != "" {
		os.Setenv("PUBLIC_API_DOMAIN", originalApiDomain)
	} else {
		os.Unsetenv("PUBLIC_API_DOMAIN")
	}
}
