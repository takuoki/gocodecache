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
  "fmt"

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

  vLevel := cache.MustGetValue(ctx, "visibility_level", "2")
}
```

## Datasource

### Raw

Define a datasource as Golang code.

```go
datasource := cache.RawSource(map[[cache.MaxKeyLength]string]string{
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

### PostgreSQL

It will be implemented soon...
