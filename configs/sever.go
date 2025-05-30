package configs

import (
	"fmt"
	"pb_launcher/utils/networktools"

	"github.com/pocketbase/pocketbase/apis"
)

func NewPBServeConfig(c Config) (*apis.ServeConfig, error) {
	ip, port, err := networktools.GetAvailablePort(c.GetBindAddress())
	if err != nil {
		return nil, err
	}
	return &apis.ServeConfig{
		ShowStartBanner: false,
		HttpAddr:        fmt.Sprintf("%s:%d", ip, port),
		AllowedOrigins:  []string{"*"},
	}, nil
}
