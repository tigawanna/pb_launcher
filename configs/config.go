package configs

import (
	"log/slog"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase/apis"
)

type Configs struct {
	*apis.ServeConfig
	ReleaseSyncInterval time.Duration  // default: 10m
	GithubRepository    string         // default: pocketbase/pocketbase
	ReleaseFilePattern  *regexp.Regexp // default: pocketbase_.+_linux_amd64\.zip
	DownloadDir         string         // default: ./downloads
}

func ReadConfigs() *Configs {
	httpAddr, ok := os.LookupEnv("ADDRESS")
	if !ok {
		httpAddr = "0.0.0.0:7090"
	}

	syncInterval := 10 * time.Minute
	repoName := "pocketbase/pocketbase"
	rgxPattern := `pocketbase_.+_linux_amd64\.zip`
	downloadDir := "./downloads"

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

	if envRepo, ok := os.LookupEnv("GITHUB_REPOSITORY"); ok {
		repoName = strings.Trim(strings.TrimSpace(envRepo), "/")
	}

	if envPattern, ok := os.LookupEnv("RELEASE_FILE_PATTERN"); ok {
		if _, err := regexp.Compile(envPattern); err == nil {
			rgxPattern = envPattern
		}
	}

	if envDir, ok := os.LookupEnv("DOWNLOAD_DIR"); ok {
		downloadDir = strings.TrimSpace(envDir)
	}

	compiledRgx, _ := regexp.Compile(rgxPattern)

	return &Configs{
		ServeConfig: &apis.ServeConfig{
			ShowStartBanner: true,
			HttpAddr:        httpAddr,
			AllowedOrigins:  []string{"*"},
		},
		ReleaseSyncInterval: syncInterval,
		GithubRepository:    repoName,
		ReleaseFilePattern:  compiledRgx,
		DownloadDir:         downloadDir,
	}
}
