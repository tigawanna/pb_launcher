package configs

import (
	"errors"
	"log/slog"
	"net"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config interface {
	GetReleaseSyncInterval() time.Duration
	GetDownloadDir() string
	GetDataDir() string
	GetPublicApiDomain() string
	GetBindAddress() string
	GetBindPort() string
}

type configs struct {
	ReleaseSyncInterval string `mapstructure:"release_sync_interval"` // default: 10m
	DownloadDir         string `mapstructure:"download_dir"`          // default: ./downloads
	DataDir             string `mapstructure:"data_dir"`              // default: ./data
	PublicApiDomain     string `mapstructure:"public_api_domain"`
	BindAddress         string `mapstructure:"bind_address"` // default: 127.0.0.1
	BindPort            string `mapstructure:"bind_port"`    // default: 8072
}

var _ Config = (*configs)(nil)

const min_sync_interval = 5 * time.Minute

func (c *configs) GetReleaseSyncInterval() time.Duration {
	duration, err := time.ParseDuration(c.ReleaseSyncInterval)
	if err != nil {
		slog.Warn("Failed to parse release sync interval",
			slog.String("raw_value", c.ReleaseSyncInterval),
			slog.String("error", err.Error()),
		)
	}
	return max(duration, min_sync_interval)
}
func (c *configs) GetDownloadDir() string {
	if c.DownloadDir == "" {
		return "./downloads"
	}
	return c.DownloadDir
}

func (c *configs) GetDataDir() string {
	if c.DataDir == "" {
		return "./data"
	}
	return c.DataDir
}

func (c *configs) GetPublicApiDomain() string {
	if c.PublicApiDomain == "" {
		return "pb.labenv.test:7080"
	}
	return c.PublicApiDomain
}
func (c *configs) GetBindAddress() string {
	if c.BindAddress == "" {
		return "127.0.0.1"
	}
	return c.BindAddress
}

func (c *configs) GetBindPort() string {
	if c.BindPort == "" {
		return "7080"
	}
	return c.BindPort
}

func loadConfigFromFile(filePath string) (*configs, error) {
	v := viper.New()
	v.SetConfigFile(filePath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		slog.Error("failed to read config file", "file", path.Base(filePath), "error", err)
		return nil, err
	}

	var cfg configs
	if err := v.Unmarshal(&cfg); err != nil {
		slog.Error("failed to unmarshal config", "file", path.Base(filePath), "error", err)
		return nil, err
	}

	return &cfg, nil
}

func LoadConfigs(configPath string) (Config, error) {
	if configPath == "" {
		return &configs{}, nil
	}
	c, err := loadConfigFromFile(configPath)
	if err != nil {
		return nil, err
	}

	c.ReleaseSyncInterval = strings.TrimSpace(c.ReleaseSyncInterval)
	c.DownloadDir = strings.TrimSpace(c.DownloadDir)
	c.DataDir = strings.TrimSpace(c.DataDir)
	c.PublicApiDomain = strings.TrimSpace(c.PublicApiDomain)
	c.BindAddress = strings.TrimSpace(c.BindAddress)
	c.BindPort = strings.TrimSpace(c.BindPort)

	if c.ReleaseSyncInterval != "" {
		duration, err := time.ParseDuration(c.ReleaseSyncInterval)
		if err != nil {
			slog.Warn("Invalid release_sync_interval format",
				slog.String("value", c.ReleaseSyncInterval),
				slog.String("error", err.Error()),
				slog.String("using_default", min_sync_interval.String()))

			c.ReleaseSyncInterval = min_sync_interval.String()
		} else if duration < min_sync_interval {
			slog.Warn("Configured release_sync_interval is too short",
				slog.Duration("provided", duration),
				slog.Duration("minimum_allowed", min_sync_interval))

			c.ReleaseSyncInterval = min_sync_interval.String()
		} else {
			c.ReleaseSyncInterval = duration.String()
		}
	}

	if net.ParseIP(c.GetBindAddress()) == nil {
		slog.Error("Invalid bind_address: not a valid IP address",
			slog.String("value", c.GetBindAddress()))
		return nil, errors.New("invalid bind_address")
	}

	portNum, err := strconv.Atoi(c.GetBindPort())
	if err != nil || portNum < 1 || portNum > 65535 {
		return nil, errors.New("invalid bind_port: must be an integer between 1 and 65535")
	}
	slog.Info("Loaded config file", slog.String("file_path", configPath))
	return c, nil
}
