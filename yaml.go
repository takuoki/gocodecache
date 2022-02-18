package gocodecache

import (
	"context"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type yamlSource struct {
	filepath string
	rootKey  string
}

func YAMLSource(filepath, rootKey string) Datasource {
	return &yamlSource{
		filepath: filepath,
		rootKey:  rootKey,
	}
}

func (d *yamlSource) ReadAll(ctx context.Context, keyLength int) (map[[MaxKeyLength]string]string, error) {
	return d.ReadFirstKeys(ctx, keyLength, nil)
}

func (d *yamlSource) ReadFirstKeys(ctx context.Context, keyLength int, firstKeys map[string]struct{}) (map[[MaxKeyLength]string]string, error) {
	f, err := os.Open(d.filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return d.read(ctx, keyLength, firstKeys, f)
}

func (d *yamlSource) read(ctx context.Context, keyLength int, firstKeys map[string]struct{}, r io.Reader) (map[[MaxKeyLength]string]string, error) {
	var m map[string]interface{}
	if err := yaml.NewDecoder(r).Decode(&m); err != nil {
		return nil, fmt.Errorf("failed to decode file: %w", err)
	}

	for k, v := range m {
		if k == d.rootKey {
			if codes, ok := v.(map[interface{}]interface{}); ok {
				return convert(codes, keyLength, firstKeys)
			} else {
				return nil, fmt.Errorf("root type must be map[interface{}]interface{} (type: %T)", v)
			}
		}
	}

	return nil, fmt.Errorf("'%s' field is not defined in the file", d.rootKey)
}
