package gocodecache_test

import (
	"context"
	"fmt"

	cache "github.com/takuoki/gocodecache"
)

func Example() {
	ctx := context.Background()

	if err := cache.InitializeGlobalCache(ctx, cache.YAMLSource("./sample/codes.yaml", "CODES"), 2); err != nil {
		// handle error
	}

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

	c, err := cache.New(ctx, cache.YAMLSource("./sample/codes_lang.yaml", "CODES"), 3)
	if err != nil {
		// handle error
	}

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
