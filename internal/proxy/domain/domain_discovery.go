package domain

import (
	"context"
	"errors"
	"log/slog"
	"pb_launcher/internal/proxy/domain/repositories"
	"time"

	"github.com/allegro/bigcache/v3"
)

type DomainServiceDiscovery struct {
	repo  repositories.ServiceRepository
	cache *bigcache.BigCache
}

func NewDomainServiceDiscovery(repo repositories.ServiceRepository) (*DomainServiceDiscovery, error) {

	cache, err := bigcache.New(context.Background(), bigcache.Config{
		Shards:           256,              // increases parallelism
		LifeWindow:       15 * time.Minute, // cache entries live longer
		CleanWindow:      30 * time.Minute, // less frequent cleanup
		MaxEntrySize:     512,              // supports moderately sized payloads
		Verbose:          false,
		HardMaxCacheSize: 128, // ~128 MB max cache size
		StatsEnabled:     false,
	})
	if err != nil {
		return nil, err
	}

	return &DomainServiceDiscovery{
		repo:  repo,
		cache: cache,
	}, nil
}

func (s *DomainServiceDiscovery) FindServiceIDByDomain(ctx context.Context, domain string) (*string, error) {
	if data, err := s.cache.Get(domain); err == nil {
		cachedID := string(data)
		if cachedID != "" {
			return &cachedID, nil
		}
	} else if !errors.Is(err, bigcache.ErrEntryNotFound) {
		slog.Warn("failed to access cache", "domain", domain, "error", err)
	}

	serviceID, err := s.repo.FindServiceIDByDomain(ctx, domain)
	if err != nil {
		return nil, err
	}

	if serviceID != nil && *serviceID != "" {
		if err := s.cache.Set(domain, []byte(*serviceID)); err != nil {
			slog.Warn("failed to cache service ID", "domain", domain, "error", err)
		}
	}

	return serviceID, nil
}

func (s *DomainServiceDiscovery) InvalidateDomain(domain string) {
	err := s.cache.Delete(domain)
	if err != nil && !errors.Is(err, bigcache.ErrEntryNotFound) {
		slog.Error("failed to invalidate domain cache", "domain", domain, "error", err)
	}
}
