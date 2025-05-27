package configs

import (
	"errors"
	"log/slog"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Configs struct {
	ReleaseSyncInterval time.Duration // default: 10m
	DownloadDir         string        // default: ./downloads
	DataDir             string        // default: ./data
	PublicApiDomain     string
	BindAddress         string // default: 127.0.0.1
	BindPort            string // default: 8072
}

func ReadConfigs() (*Configs, error) {
	syncInterval := 10 * time.Minute
	downloadDir := "./downloads"
	dataDir := "./data"
	bindAddress := "127.0.0.1"
	bindPort := "8072"

	const minSyncInterval = 5 * time.Minute
	if envInterval, ok := os.LookupEnv("RELEASE_SYNC_INTERVAL"); ok {
		duration, err := time.ParseDuration(envInterval)
		if err != nil {
			slog.Warn("invalid RELEASE_SYNC_INTERVAL format, using default", "error", err, "default", minSyncInterval)
			syncInterval = minSyncInterval
		} else if duration < minSyncInterval {
			slog.Warn("configured RELEASE_SYNC_INTERVAL is too short, using minimum allowed", "provided", duration, "minimum", minSyncInterval)
			syncInterval = minSyncInterval
		} else {
			syncInterval = duration
		}
	}

	if envDir, ok := os.LookupEnv("DOWNLOAD_DIR"); ok {
		downloadDir = strings.TrimSpace(envDir)
	}
	if envDir, ok := os.LookupEnv("DATA_DIR"); ok {
		dataDir = strings.TrimSpace(envDir)
	}

	if envAddr, ok := os.LookupEnv("BIND_ADDRESS"); ok {
		bindAddress = strings.TrimSpace(envAddr)
	}
	if net.ParseIP(bindAddress) == nil {
		return nil, errors.New("invalid BIND_ADDRESS: not a valid IP address")
	}

	if envPort, ok := os.LookupEnv("BIND_PORT"); ok {
		bindPort = strings.TrimSpace(envPort)
	}
	portNum, err := strconv.Atoi(bindPort)
	if err != nil || portNum < 1 || portNum > 65535 {
		return nil, errors.New("invalid BIND_PORT: must be an integer between 1 and 65535")
	}

	apiDomain, ok := os.LookupEnv("PUBLIC_API_DOMAIN")
	if !ok || strings.TrimSpace(apiDomain) == "" {
		return nil, errors.New("missing required environment variable: PUBLIC_API_DOMAIN")
	}

	return &Configs{
		ReleaseSyncInterval: syncInterval,
		DownloadDir:         downloadDir,
		DataDir:             dataDir,
		PublicApiDomain:     apiDomain,
		BindAddress:         bindAddress,
		BindPort:            bindPort,
	}, nil
}
