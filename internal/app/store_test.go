package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const backFile = "tempfile"

func Test_mapStorage_store(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		useBackFile bool
	}{
		{
			name: "simple url",
			args: args{
				url: "http://www.kiae.su/",
			},
			wantErr: false,
		},
		{
			name: "with back store",
			args: args{
				url: "http://www.kiae.su/",
			},
			wantErr:     false,
			useBackFile: true,
		},
	}
	for _, tt := range tests {
		var fileName string
		if tt.useBackFile {
			fileName = backFile
		}
		t.Run(tt.name, func(t *testing.T) {
			sh := newStoreWithFile(fileName)
			newKey, err := sh.s.store(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("storage.store() error = %v, wantErr %v", err, tt.wantErr)
			}
			url, err := sh.s.get(newKey)
			if err != nil {
				t.Errorf("get returns error %s while should not", err)
			}
			if url != tt.args.url {
				t.Errorf("stored url does not match original %s", tt.args.url)
			}
			if tt.useBackFile {
				sh.s.close()
				assert.FileExists(t, backFile)
				cont, err := os.ReadFile(backFile)
				assert.NoError(t, err)
				assert.Contains(t, string(cont), newKey)
				assert.Contains(t, string(cont), tt.args.url)
			}
		})
		if tt.useBackFile {
			os.Remove(backFile)
		}
	}
}

func Test_newStoreWithFile(t *testing.T) {
	tests := []struct {
		name            string
		useBackFile     bool
		backFileContent string
		wantKey         string
		wantVal         string
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name:    "Simple",
			wantErr: assert.NoError,
		},
		{
			name:            "With back store",
			wantErr:         assert.NoError,
			useBackFile:     true,
			backFileContent: "abc def",
			wantKey:         "abc",
			wantVal:         "def",
		},
	}
	for _, tt := range tests {
		if tt.useBackFile {
			os.Remove(backFile)
		}
		if tt.backFileContent != "" {
			f, err := os.Create(backFile)
			require.NoError(t, err)
			f.Write([]byte(tt.backFileContent))
			f.Close()
		}
		t.Run(tt.name, func(t *testing.T) {
			var fileName string
			if tt.useBackFile {
				fileName = backFile
			}
			st := newStoreWithFile(fileName)
			tt.wantErr(t, st.s.error())
			if tt.wantKey != "" {
				val, err := st.s.get(tt.wantKey)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantVal, val)
			}
		})
		if tt.useBackFile {
			os.Remove(backFile)
		}
	}
}
