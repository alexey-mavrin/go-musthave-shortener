package app

import (
	"encoding/json"
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

func Test_storeJSONHandler(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "simple url",
			body: `{"url":"http://www.kiae.su"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sh := newStore()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost,
				"/api/shorten", strings.NewReader(tt.body))
			sh.storeJSONHandler(w, r)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, http.StatusCreated, w.Code)

			body, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			short := new(Result)
			err = json.Unmarshal(body, short)
			assert.NoError(t, err)
			assert.NotEqual(t, "", short.Result)
			assert.True(t, strings.HasPrefix(short.Result, serverAddress),
				"response contains server address")
			assert.Greater(t, len(short.Result), len(serverAddress))
			assert.Equal(t, w.Header().Get("Content-Type"), "application/json")
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
