package gocodecache

import (
	"context"
	"fmt"
	"strings"
)

type Datasource interface {
	ReadAll(ctx context.Context, keyLength int) (Codes, error)
	ReadFirstKeys(ctx context.Context, keyLength int, firstKeys map[string]struct{}) (Codes, error)
}

func convert(m map[interface{}]interface{}, keyLength int, firstKeys map[string]struct{}) (Codes, error) {
	result := Codes{}
	if err := convertKeys(m, keyLength, 0, nil, firstKeys, result); err != nil {
		return nil, err
	}
	return result, nil
}

func convertKeys(m map[interface{}]interface{}, max, index int, keys []string, firstKeys map[string]struct{}, result Codes) error {
	if index < max-1 {
		for k, v := range m {
			key := fmt.Sprint(k)
			if index == 0 && firstKeys != nil {
				if _, ok := firstKeys[key]; !ok {
					continue
				}
			}
			keys := append(keys, key)
			if v2, ok := v.(map[interface{}]interface{}); ok {
				if err := convertKeys(v2, max, index+1, keys, firstKeys, result); err != nil {
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
