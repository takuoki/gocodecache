package gocodecache

import (
	"context"
	"io"
)

func ReadForYAML(ctx context.Context, rootKey string, keyLength int, firstKeys map[string]struct{}, r io.Reader) (map[[MaxKeyLength]string]string, error) {
	return (&yamlSource{rootKey: rootKey}).read(ctx, keyLength, firstKeys, r)
}
