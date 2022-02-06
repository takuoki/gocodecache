package gocodecache_test

import (
	"context"
	"fmt"
	"time"

	cache "github.com/takuoki/gocodecache"
)

func Example() {
	ctx := context.Background()

	// -- How to initialize --

	if err := cache.InitializeGlobalCache(ctx, cache.YAMLSource("./sample/codes.yaml", "CODES"), 2); err != nil {
		// handle error
	}

	// -- How to get value --

	accType1Str, err := cache.GetValue(ctx, "account_type", "1")
	if err != nil {
		// handle error
	}

	fmt.Println(accType1Str)
	fmt.Println(cache.MustGetValue(ctx, "visibility_level", "2"))

	// Output:
	// Anonymous account
	// Public
}

func ExampleCache() {
	ctx := context.Background()

	// -- How to initialize --

	c, err := cache.New(ctx, cache.YAMLSource("./sample/codes_lang.yaml", "CODES"), 3)
	if err != nil {
		// handle error
	}

	// -- How to get value --

	accType1jaStr, err := c.GetValue(ctx, "account_type", "1", "ja-JP")
	if err != nil {
		// handle error
	}

	fmt.Println(accType1jaStr)
	fmt.Println(c.MustGetValue(ctx, "visibility_level", "2", "en-US"))

	// Output:
	// 匿名アカウント
	// Public
}

func ExamplePostgresSource() {
	ctx := context.Background()

	// -- How to initialize --

	db, err := cache.ConnectPostgres(
		"localhost",
		"5432",
		"root",
		"root",
		"postgres",
		"disable",
	)
	if err != nil {
		// handle error
	}
	defer db.Close()

	for i := 0; ; i++ {
		if i >= dbPingRetryLimit {
			// handle error
		}
		if err = db.Ping(); err == nil {
			break
		}
		time.Sleep(dbPingRetryInterval)
	}

	c, err := cache.New(ctx, cache.PostgresSource(db, "codes", [cache.MaxKeyLength]string{"key1", "key2"}, "value"), 2)
	if err != nil {
		// handle error
	}

	// -- How to get value --

	accType1jaStr, err := c.GetValue(ctx, "account_type", "1")
	if err != nil {
		// handle error
	}

	fmt.Println(accType1jaStr)
	fmt.Println(c.MustGetValue(ctx, "visibility_level", "2"))

	// Output:
	// Anonymous account
	// Public
}
