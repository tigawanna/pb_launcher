package configs

import (
	"errors"
	"log/slog"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type TlsConfig interface {
	GetProvider() string
	GetProp(key string) (string, bool)
}

type Config interface {
	GetReleaseSyncInterval() time.Duration
	GetCommandCheckInterval() time.Duration
	GetCertificateCheckInterval() time.Duration

	GetDownloadDir() string
	GetDataDir() string
	GetCertificatesDir() string
	GetMinCertificateTtl() time.Duration

	GetDomain() string
	GetBindAddress() string
	GetBindPort() string

	IsHttpsEnabled() bool
	IsHttpsRedirectDisabled() bool

	GetBindHttpsPort() string
	GetTlsConfig() TlsConfig
}

type tls_configs struct {
	Provider string            `mapstructure:"provider"`
	Props    map[string]string `mapstructure:"props"`
}

var _ TlsConfig = (*tls_configs)(nil)

func (c *tls_configs) GetProvider() string {
	provider := strings.TrimSpace(c.Provider)
	if provider == "" {
		slog.Warn("TLS provider is empty, using default 'selfsigned'")
		return "selfsigned"
	}
	return provider
}

func (c *tls_configs) GetProp(key string) (string, bool) {
	if c.Props == nil {
		return "", false
	}
	val, ok := c.Props[key]
	return val, ok
}

type configs struct {
	ReleaseSyncInterval      string `mapstructure:"release_sync_interval"`      // default: 10m
	CommandCheckInterval     string `mapstructure:"command_check_interval"`     // default: 10ms
	CertificateCheckInterval string `mapstructure:"certificate_check_interval"` // default: 1h

	DownloadDir          string `mapstructure:"download_dir"`     // default: ./downloads
	CertificatesDir      string `mapstructure:"certificates_dir"` // default: ./.certificates
	DataDir              string `mapstructure:"data_dir"`         // default: ./data
	Domain               string `mapstructure:"domain"`
	BindAddress          string `mapstructure:"bind_address"` // default: 127.0.0.1
	BindPort             string `mapstructure:"bind_port"`    // default: 8072
	Https                bool   `mapstructure:"https"`
	DisableHttpsRedirect bool   `mapstructure:"disable_https_redirect"`
	HttpsPort            string `mapstructure:"https_port"`          // default: 8443
	MinCertificateTtl    string `mapstructure:"min_certificate_ttl"` // default: 720h

	Tls tls_configs `mapstructure:"cert"`
}

var _ Config = (*configs)(nil)

const min_sync_interval = 5 * time.Minute
const min_command_check_interval = 10 * time.Second
const min_certificate_check_interval = time.Minute
const min_certificate_ttl = 30 * 24 * time.Hour

func (c *configs) GetReleaseSyncInterval() time.Duration {
	duration, err := time.ParseDuration(c.ReleaseSyncInterval)
	if err != nil {
		slog.Warn("Failed to parse release_sync_interval",
			slog.String("raw_value", c.ReleaseSyncInterval),
			slog.String("error", err.Error()),
		)
		return min_sync_interval
	}
	return max(duration, min_sync_interval)
}

func (c *configs) GetCommandCheckInterval() time.Duration {
	duration, err := time.ParseDuration(c.CommandCheckInterval)
	if err != nil {
		slog.Warn("Failed to parse command_check_interval",
			slog.String("raw_value", c.CommandCheckInterval),
			slog.String("error", err.Error()),
		)
		return min_command_check_interval
	}
	return max(duration, min_command_check_interval)
}
func (c *configs) GetCertificateCheckInterval() time.Duration {
	duration, err := time.ParseDuration(c.CertificateCheckInterval)
	if err != nil {
		slog.Warn("Failed to parse certificate_check_interval",
			slog.String("raw_value", c.CertificateCheckInterval),
			slog.String("error", err.Error()),
		)
		return min_certificate_check_interval
	}
	return max(duration, min_certificate_check_interval)
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

func (c *configs) GetCertificatesDir() string {
	if c.CertificatesDir == "" {
		return "./.certificates"
	}
	return c.CertificatesDir
}

func (c *configs) GetDomain() string {
	if c.Domain == "" {
		return "pb.labenv.test"
	}
	return strings.Split(c.Domain, ":")[0]
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

func (c *configs) IsHttpsEnabled() bool          { return c.Https }
func (c *configs) IsHttpsRedirectDisabled() bool { return c.DisableHttpsRedirect }

func (c *configs) GetBindHttpsPort() string {
	if c.HttpsPort == "" {
		return "8443"
	}
	return c.HttpsPort
}

func (c *configs) GetMinCertificateTtl() time.Duration {
	duration, err := time.ParseDuration(c.MinCertificateTtl)
	if err != nil {
		slog.Warn("Failed to parse min_certificate_ttl",
			slog.String("raw_value", c.MinCertificateTtl),
			slog.String("error", err.Error()),
		)
		return min_certificate_ttl
	}
	return max(duration, min_certificate_ttl)
}

func (c *configs) GetTlsConfig() TlsConfig { return &c.Tls }

func loadConfigFromFile(filePath string) (*configs, error) {
	v := viper.New()
	v.SetConfigFile(filePath)
	v.SetConfigType("yaml")

	var cfg configs

	if err := v.ReadInConfig(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &cfg, nil
		}
		slog.Error("failed to read config file", "file", path.Base(filePath), "error", err)
		return nil, err
	}

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
	c.Domain = strings.TrimSpace(c.Domain)
	c.BindAddress = strings.TrimSpace(c.BindAddress)
	c.BindPort = strings.TrimSpace(c.BindPort)
	c.HttpsPort = strings.TrimSpace(c.HttpsPort)

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

	if c.CommandCheckInterval != "" {
		duration, err := time.ParseDuration(c.CommandCheckInterval)
		if err != nil {
			slog.Warn("Invalid command_check_interval format",
				slog.String("value", c.CommandCheckInterval),
				slog.String("error", err.Error()),
				slog.String("using_default", min_command_check_interval.String()))

			c.CommandCheckInterval = min_command_check_interval.String()
		} else if duration < min_command_check_interval {
			slog.Warn("Configured command_check_interval is too short",
				slog.Duration("provided", duration),
				slog.Duration("minimum_allowed", min_command_check_interval))

			c.CommandCheckInterval = min_command_check_interval.String()
		} else {
			c.CommandCheckInterval = duration.String()
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
