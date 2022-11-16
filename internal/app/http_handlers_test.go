package app

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_storeHandler_storeHandler(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "simple url",
			url:  "http://www.kiae.su/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sh := newStore()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost,
				"/", strings.NewReader(tt.url))
			sh.storeHandler(w, r)
			res := w.Result()
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusCreated, w.Code)
			assert.True(t, strings.HasPrefix(string(body), serverAddress),
				"response contains server address")
			assert.Greater(t, len(string(body)), len(serverAddress))
		})
	}
}

func Test_storeHandler_fetchHandler(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "simple url",
			url:  "http://www.kiae.su",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sh := newStore()
			key, err := sh.s.store(tt.url)
			assert.NoError(t, err)

			r := newServer(&sh)
			ts := httptest.NewServer(r)
			defer ts.Close()

			req, err := http.NewRequest(http.MethodGet, ts.URL+"/"+key, nil)
			assert.NoError(t, err)
			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}
			res, err := client.Do(req)
			assert.NoError(t, err)

			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
			assert.Contains(t, string(body),
				`<a href="`+tt.url+`">Temporary Redirect</a>.`)
		})
	}
}
