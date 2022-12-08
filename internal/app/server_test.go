package app

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newServer(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		method       string
		reqBody      []byte
		resBody      []byte
		resBodyRegex string
	}{
		{
			name:    "basic",
			path:    "/",
			method:  http.MethodGet,
			reqBody: nil,
			resBody: []byte{},
		},
		{
			name:         "plain",
			path:         "/",
			method:       http.MethodPost,
			reqBody:      []byte("http://www.kiae.su"),
			resBodyRegex: `^http:\/\/localhost:8080\/[a-zA-Z]{6}`,
		},
		{
			name:         "json",
			path:         "/api/shorten",
			method:       http.MethodPost,
			reqBody:      []byte(`{"url":"http://www.kiae.su"}`),
			resBodyRegex: `^{"result":"http:\/\/localhost:8080\/[a-zA-Z]{6}"}`,
		},
	}

	c := Config{
		BaseURL: "http://localhost:8080/",
	}
	c.sh = newStore()
	r := newServer(c)
	srv := httptest.NewServer(r)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method,
				srv.URL+tt.path,
				bytes.NewReader(tt.reqBody),
			)
			require.NoError(t, err)
			res, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			t.Log("body: ", string(body))

			require.NoError(t, err)
			if tt.resBody != nil {
				assert.Equal(t, body, tt.resBody)
			} else {
				assert.Regexp(t, regexp.MustCompile(tt.resBodyRegex), string(body))
			}
		})
	}
}
