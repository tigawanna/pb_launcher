package domain

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"log/slog"
	"pb_launcher/internal/proxy/domain/dtos"
	"pb_launcher/internal/proxy/domain/repositories"

	"time"

	"github.com/allegro/bigcache/v3"
)

type ProxyEntryDiscovery struct {
	repo  repositories.ProxyEntriesRepository
	cache *bigcache.BigCache
}

func init() {
	gob.Register(&dtos.ProxyEntryDto{})
}

func NewProxyEntryDiscovery(repo repositories.ProxyEntriesRepository) (*ProxyEntryDiscovery, error) {

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

	return &ProxyEntryDiscovery{
		repo:  repo,
		cache: cache,
	}, nil
}

func (s *ProxyEntryDiscovery) FindEnabledProxyEntryByID(ctx context.Context, id string) (*dtos.ProxyEntryDto, error) {
	if data, err := s.cache.Get(id); err == nil {
		buf := bytes.NewBuffer(data)
		dec := gob.NewDecoder(buf)
		var dto dtos.ProxyEntryDto
		if err := dec.Decode(&dto); err == nil {
			return &dto, nil
		}
	} else if !errors.Is(err, bigcache.ErrEntryNotFound) {
		slog.Warn("failed to access cache", "proxy_entry", id, "error", err)
	}

	dto, err := s.repo.FindEnabledProxyEntryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(dto); err == nil {
		if err := s.cache.Set(id, buf.Bytes()); err != nil {
			slog.Warn("failed to cache proxy entries", "proxy_entry", id, "error", err)
		}
	}

	return dto, nil
}

func (s *ProxyEntryDiscovery) InvalidateProxyEntriesCacheByID(id string) error {
	err := s.cache.Delete(id)
	if err != nil && !errors.Is(err, bigcache.ErrEntryNotFound) {
		slog.Error("failed to invalidate cache", "id", id, "error", err)
		return err
	}
	slog.Info("invalidated proxy entry cache", "proxy_entry", id)
	return nil
}
