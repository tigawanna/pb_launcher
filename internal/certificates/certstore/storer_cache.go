package certstore

import (
	"bytes"
	"context"
	"encoding/gob"
	"pb_launcher/internal/certificates/tlscommon"
	"time"

	"github.com/allegro/bigcache/v3"
)

func init() {
	gob.Register(tlscommon.Certificate{})
}

type TlsStorerCache struct {
	storer *TlsStorer
	cache  *bigcache.BigCache
}

var _ tlscommon.Store = (*TlsStorerCache)(nil)

func NewTlsStorerCache(storer *TlsStorer) (*TlsStorerCache, error) {
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

	return &TlsStorerCache{
		storer: storer,
		cache:  cache,
	}, nil
}

func (t *TlsStorerCache) Store(domain string, cert tlscommon.Certificate) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(cert); err != nil {
		return err
	}

	if err := t.cache.Set(domain, buf.Bytes()); err != nil {
		return err
	}

	return t.storer.Store(domain, cert)
}

func (t *TlsStorerCache) Resolve(domain string) (*tlscommon.Certificate, error) {
	data, err := t.cache.Get(domain)
	if err == nil {
		var cert tlscommon.Certificate
		if decodeErr := gob.NewDecoder(bytes.NewReader(data)).
			Decode(&cert); decodeErr == nil {
			return &cert, nil
		}
	}

	cert, err := t.storer.Resolve(domain)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if encodeErr := gob.NewEncoder(&buf).Encode(cert); encodeErr == nil {
		_ = t.cache.Set(domain, buf.Bytes()) // ignore error
	}

	return cert, nil
}
