package gocodecache

import (
	"context"
	"fmt"
	"strings"
)

type Datasource interface {
	ReadAll(ctx context.Context, keyLength int) (map[[MaxKeyLength]string]string, error)
}

func convert(m map[interface{}]interface{}, keyLength int) (map[[MaxKeyLength]string]string, error) {
	result := map[[MaxKeyLength]string]string{}
	if err := convertKeys(m, keyLength, 0, nil, result); err != nil {
		return nil, err
	}
	return result, nil
}

func convertKeys(m map[interface{}]interface{}, max, index int, keys []string, result map[[MaxKeyLength]string]string) error {
	if index < max-1 {
		for k, v := range m {
			keys := append(keys, fmt.Sprint(k))
			if v2, ok := v.(map[interface{}]interface{}); ok {
				if err := convertKeys(v2, max, index+1, keys, result); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("value type must be map[interface{}]interface{} (key: [%s], type: %T)", strings.Join(keys, ", "), v)
			}
		}
	} else if index >= max-1 {
		for k, v := range m {
			keys := append(keys, fmt.Sprint(k))
			if v2, ok := v.(string); ok {
				keys2 := [MaxKeyLength]string{}
				for i, key := range keys {
					keys2[i] = key
				}
				result[keys2] = v2
			} else {
				return fmt.Errorf("value type must be string (key: [%s], type: %T)", strings.Join(keys, ", "), v)
			}
		}
	}

	return nil
}
