package gocodecache

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

const MaxKeyLength = 5

var (
	ErrNotInitialized = errors.New("cache is not initialized")
	ErrCodeNotFound   = errors.New("code not found")
)

type Cache struct {
	mu        sync.RWMutex
	codes     map[[MaxKeyLength]string]string
	ds        Datasource
	keyLength int
}

var defaultCache *Cache

func InitializeGlobalCache(ctx context.Context, ds Datasource, keyLength int) error {
	c, err := New(ctx, ds, keyLength)
	if err != nil {
		return err
	}

	defaultCache = c

	return nil
}

func Reload(ctx context.Context) error {
	if defaultCache == nil {
		return ErrNotInitialized
	}
	return defaultCache.Reload(ctx)
}

func GetValue(ctx context.Context, keys ...string) (string, error) {
	if defaultCache == nil {
		return "", ErrNotInitialized
	}
	return defaultCache.GetValue(ctx, keys...)
}

func MustGetValue(ctx context.Context, keys ...string) string {
	if defaultCache == nil {
		panic(ErrNotInitialized)
	}
	return defaultCache.MustGetValue(ctx, keys...)
}

func New(ctx context.Context, ds Datasource, keyLength int) (*Cache, error) {
	if ds == nil {
		return nil, errors.New("datasource is nil")
	}
	if keyLength < 1 || MaxKeyLength < keyLength {
		return nil, fmt.Errorf("invalid key length, must be between 1 and %d", MaxKeyLength)
	}

	c := &Cache{
		ds:        ds,
		keyLength: keyLength,
	}
	if err := c.load(ctx); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Cache) load(ctx context.Context) error {
	m, err := c.ds.ReadAll(ctx, c.keyLength)
	if err != nil {
		return fmt.Errorf("failed to read from datasource: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.codes = m

	return nil
}

func (c *Cache) Reload(ctx context.Context) error {
	return c.load(ctx)
}

func (c *Cache) GetValue(ctx context.Context, keys ...string) (string, error) {
	k := [MaxKeyLength]string{}
	for i, key := range keys {
		k[i] = key
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.codes[k]
	if !ok {
		return "", ErrCodeNotFound
	}

	return v, nil
}

func (c *Cache) MustGetValue(ctx context.Context, keys ...string) string {
	v, err := c.GetValue(ctx, keys...)
	if err != nil {
		panic(err)
	}
	return v
}