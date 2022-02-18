# gocodecache

[![GoDoc](https://godoc.org/github.com/takuoki/gocodecache?status.svg)](https://godoc.org/github.com/takuoki/gocodecache)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
![CI](https://github.com/takuoki/gocodecache/actions/workflows/auto-test.yml/badge.svg)

An in-memory cache library for code master in Golang.

## Installation

```
go get github.com/takuoki/gocodecache
```

## Get started

### Create YAML file

```yaml
# ./sample/codes.yaml
VERSION: 0.1.0

CODES:
  account_type:
    1: Anonymous account
    2: General account
    3: Administrator account
  visibility_level:
    1: Private
    2: Public
```

### Code

```go
package main

import (
  "context"

  cache "github.com/takuoki/gocodecache"
)

func main() {
  ctx := context.Background()

  if err := cache.InitializeGlobalCache(ctx, cache.YAMLSource("./sample/codes.yaml", "CODES"), 2); err != nil {
    // handle error
  }

  accType1Str, err := cache.GetValue(ctx, "account_type", "1")
  if err != nil {
    // handle error
  }

  vLevel2Str := cache.MustGetValue(ctx, "visibility_level", "2")
}
```

## Datasource

### Raw

Define a datasource as Golang code.

```go
datasource := cache.RawSource(cache.Codes{
  {"account_type", "1"}:     "Anonymous account",
  {"account_type", "2"}:     "General account",
  {"account_type", "3"}:     "Administrator account",
  {"visibility_level", "1"}: "Private",
  {"visibility_level", "2"}: "Public",
})
```

### YAML

```go
datasource := cache.YAMLSource("./sample/codes.yaml", "CODES")
```

```yaml
VERSION: 0.1.0

CODES:
  account_type:
    1: Anonymous account
    2: General account
    3: Administrator account
  visibility_level:
    1: Private
    2: Public
```

### Relational database

!!! **WARNING** !!! Tested only with PostgreSQL.

```go
datasource := cache.RdbSource(db, "codes", []string{"key1", "key2"}, "value")
```

table: codes

| key1             | key2 | value                 |
| :--------------- | :--- | :-------------------- |
| account_type     | 1    | Anonymous account     |
| account_type     | 2    | General account       |
| account_type     | 3    | Administrator account |
| visibility_level | 1    | Private               |
| visibility_level | 2    | Public                |

## Tips

### Partial loading

By specifying first keys, you can load a part of the code master.

```go
package main

import (
  "context"

  cache "github.com/takuoki/gocodecache"
)

func main() {
  ctx := context.Background()

  // ...

  c, err := cache.New(ctx, datasource, 2, cache.WithLoadFirstKeys("account_type"))
  if err != nil {
    // handle error
  }

  // found
  accType1Str, err := c.GetValue(ctx, "account_type", "1")
  if err != nil {
    // handle error
  }

  // not found
  vLevel2Str, err := c.GetValue(ctx, "visibility_level", "2")
  if err != nil {
    // handle error
  }
}
```

### Automatic reloading

By implementing the reloading process in goroutine, automatic reloading can be achieved. ([sample](sample/main.go))

```go
package main

import (
  "context"
  "time"

  cache "github.com/takuoki/gocodecache"
)

const reloadInterval = 1 * time.Hour

func main() {

  // ...

  c, err := cache.New(ctx, datasource, 2)
  if err != nil {
    // handle error
  }
  go reload(ctx, c)

  // ...

}

func reload(ctx context.Context, c *cache.Cache) {
  for {
    select {
    case <-ctx.Done():
      return
    default:
    }

    time.Sleep(reloadInterval)
    if err := c.Reload(ctx); err != nil {
      // handle error
    }
  }
}
```

### Internationalization (I18n)

I18n can be supported by adding language codes to the keys.

```yaml
# ./sample/codes_lang.yaml
VERSION: 0.1.0

CODES:
  account_type:
    1:
      en-US: Anonymous account
      ja-JP: 匿名アカウント
    2:
      en-US: General account
      ja-JP: 一般アカウント
    3:
      en-US: Administrator account
      ja-JP: 管理者アカウント
  visibility_level:
    1:
      en-US: Private
      ja-JP: 非公開
    2:
      en-US: Public
      ja-JP: 公開
```
