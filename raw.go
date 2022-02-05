package gocodecache

import (
	"context"
)

type rawSource struct {
	data map[[MaxKeyLength]string]string
}

func RawSource(data map[[MaxKeyLength]string]string) Datasource {
	return &rawSource{data: data}
}

func (d *rawSource) ReadAll(ctx context.Context, keyLength int) (map[[MaxKeyLength]string]string, error) {
	return d.data, nil
}
