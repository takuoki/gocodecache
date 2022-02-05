package gocodecache

import (
	"context"
	"io"
)

func ReadAllForYAML(ctx context.Context, rootKey string, keyLength int, r io.Reader) (map[[MaxKeyLength]string]string, error) {
	return (&yamlSource{rootKey: rootKey}).readAll(ctx, keyLength, r)
}
