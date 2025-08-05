package repositories

import (
	"context"
	"pb_launcher/internal/proxy/domain/dtos"
)

type ProxyEntriesRepository interface {
	FindEnabledProxyEntryByID(ctx context.Context, id string) (*dtos.ProxyEntryDto, error)
}
