package gocodecache_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	cache "github.com/takuoki/gocodecache"
)

const (
	dbPingRetryLimit    = 30
	dbPingRetryInterval = 1 * time.Second
)

func TestPostgres(t *testing.T) {
	testcases := map[string]struct {
		keyLength       int
		firstKeys       []string
		tableName       string
		keyColumnNames  []string
		valueColumnName string
		want            cache.Codes
	}{
		"success: codes": {
			keyLength:       2,
			tableName:       "codes",
			keyColumnNames:  []string{"key1", "key2"},
			valueColumnName: "value",
			want:            data,
		},
		"success: codes_lang": {
			keyLength:       3,
			tableName:       "codes_lang",
			keyColumnNames:  []string{"key1", "key2", "lang"},
			valueColumnName: "value",
			want:            dataLang,
		},
		"success: len(firstKeys) == 1": {
			keyLength:       2,
			firstKeys:       []string{"account_type"},
			tableName:       "codes",
			keyColumnNames:  []string{"key1", "key2"},
			valueColumnName: "value",
			want: cache.Codes{
				{"account_type", "1"}: "Anonymous account",
				{"account_type", "2"}: "General account",
				{"account_type", "3"}: "Administrator account",
			},
		},
		"success: len(firstKeys) == 2": {
			keyLength:       2,
			firstKeys:       []string{"account_type", "visibility_level"},
			tableName:       "codes",
			keyColumnNames:  []string{"key1", "key2"},
			valueColumnName: "value",
			want:            data,
		},
	}

	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			"localhost",
			"5432",
			"root",
			"root",
			"postgres",
			"disable",
		),
	)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	for i := 0; ; i++ {
		if i >= dbPingRetryLimit {
			t.Fatalf("failed to ping (retryCount=%d): %v", i, err)
		}
		if err = db.Ping(); err == nil {
			break
		}
		time.Sleep(dbPingRetryInterval)
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			c, err := cache.New(ctx, cache.RdbSource(db, tc.tableName, tc.keyColumnNames, tc.valueColumnName),
				tc.keyLength, cache.WithLoadFirstKeys(tc.firstKeys...))
			if err != nil {
				t.Fatalf("failed to create new cache: %v", err)
			}

			for k, v := range tc.want {
				keys := []string{}
				for _, k1 := range k {
					if k1 != "" {
						keys = append(keys, k1)
					}
				}
				str, err := c.GetValue(ctx, keys...)
				if assert.Nil(t, err, "error must be nil: ", keys) {
					assert.Equal(t, v, str, "value doesn't match: ", keys)
				}
			}
		})
	}
}
