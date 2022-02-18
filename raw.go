package gocodecache

import (
	"context"
)

type rawSource struct {
	data Codes
}

func RawSource(data Codes) Datasource {
	return &rawSource{data: data}
}

func (d *rawSource) ReadAll(ctx context.Context, keyLength int) (Codes, error) {
	return d.data, nil
}

func (d *rawSource) ReadFirstKeys(ctx context.Context, keyLength int, firstKeys map[string]struct{}) (Codes, error) {
	data := Codes{}
	for k, v := range d.data {
		if _, ok := firstKeys[k[0]]; !ok {
			continue
		}
		data[k] = v
	}
	return data, nil
}
