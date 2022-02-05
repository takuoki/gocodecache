package gocodecache_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	cache "github.com/takuoki/gocodecache"
)

func TestYAMLFile(t *testing.T) {
	ctx := context.Background()

	r, err := cache.YAMLSource("./sample/codes.yaml", "CODES").ReadAll(ctx, 2)
	if err != nil {
		t.Fatalf("failed to read yaml file: %v", err)
	}

	assert.Equal(t, data, r)
}

func TestYAML(t *testing.T) {
	ctx := context.Background()

	testcases := map[string]struct {
		keyLength int
		yaml      string
		want      map[[cache.MaxKeyLength]string]string
		wantErr   string
	}{
		"success: key length == 2": {
			keyLength: 2,
			yaml: `
CODES:
  account_type:
    1: Anonymous account
    2: General account
    3: Administrator account
  visibility_level:
    1: Private
    2: Public
`,
			want: data,
		},
		"success: key length == 3": {
			keyLength: 3,
			yaml: `
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
`,
			want: dataLang,
		},
		"failure: invalid yaml": {
			keyLength: 1,
			yaml:      "{",
			wantErr:   "failed to decode file: ",
		},
		"failure: CODES not defined": {
			keyLength: 1,
			yaml:      "VERSION: 0.1.0",
			wantErr:   "'CODES' field is not defined in the file",
		},
		"failure: invalid CODES": {
			keyLength: 1,
			yaml:      "CODES: 0.1.0",
			wantErr:   `root type must be map\[interface\{\}\]interface\{\} \(type: .*\)`,
		},
		"failure: invalid key length (short)": {
			keyLength: 1,
			yaml: `
CODES:
  account_type:
    1: Anonymous account
    2: General account
    3: Administrator account
  visibility_level:
    1: Private
    2: Public
`,
			wantErr: `value type must be string \(key: \[.*\], type: .*\)`,
		},
		"failure: invalid key length (long)": {
			keyLength: 3,
			yaml: `
CODES:
  account_type:
    1: Anonymous account
    2: General account
    3: Administrator account
  visibility_level:
    1: Private
    2: Public
`,
			wantErr: `value type must be map\[interface\{\}\]interface\{\} \(key: \[.*\], type: .*\)`,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			buf := bytes.NewBufferString(tc.yaml)
			r, err := cache.ReadAllForYAML(ctx, "CODES", tc.keyLength, buf)

			if tc.wantErr == "" {
				if assert.Nil(t, err, "error must be nil") {
					assert.Equal(t, tc.want, r)
				}
			} else {
				if assert.NotNil(t, err, "error must not be nil") {
					assert.Regexp(t, "^"+tc.wantErr, err.Error())
				}
			}
		})
	}
}
