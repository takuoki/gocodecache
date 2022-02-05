package gocodecache_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	cache "github.com/takuoki/gocodecache"
)

var data = map[[cache.MaxKeyLength]string]string{
	{"account_type", "1"}:     "Anonymous account",
	{"account_type", "2"}:     "General account",
	{"account_type", "3"}:     "Administrator account",
	{"visibility_level", "1"}: "Private",
	{"visibility_level", "2"}: "Public",
}

func TestCache(t *testing.T) {

	ctx := context.Background()

	c, err := cache.New(ctx, cache.RawSource(data), 2)
	if err != nil {
		t.Fatalf("failed to create new cache: %v", err)
	}

	r1, err := c.GetValue(ctx, "account_type", "1")
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, "Anonymous account", r1)

	r2 := c.MustGetValue(ctx, "visibility_level", "2")
	assert.Equal(t, "Public", r2)

	_, err = c.GetValue(ctx, "visibility_level", "3")
	if assert.NotNil(t, err, "error must not be nil") {
		assert.Equal(t, cache.ErrCodeNotFound, err)
	}
}

func TestGlobalCache(t *testing.T) {

	ctx := context.Background()

	_, err := cache.GetValue(ctx, "account_type", "1")
	if assert.NotNil(t, err, "error must not be nil") {
		assert.Equal(t, cache.ErrNotInitialized, err)
	}

	err = cache.InitializeGlobalCache(ctx, cache.RawSource(data), 2)
	if err != nil {
		t.Fatalf("failed to create new cache: %v", err)
	}

	r1, err := cache.GetValue(ctx, "account_type", "1")
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, "Anonymous account", r1)

	r2 := cache.MustGetValue(ctx, "visibility_level", "2")
	assert.Equal(t, "Public", r2)

	_, err = cache.GetValue(ctx, "visibility_level", "3")
	if assert.NotNil(t, err, "error must not be nil") {
		assert.Equal(t, cache.ErrCodeNotFound, err)
	}
}