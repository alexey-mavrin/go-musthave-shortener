package app

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_storeHandler_dispatchHandler(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "post url to shoten",
			url:  "http://www.kiae.su",
		},
	}
	for _, tt := range tests {
		sh := newStoreHandler()
		s := httptest.NewServer(http.HandlerFunc(sh.dispatchHandler))
		defer s.Close()
		client := s.Client()

		postRes, err := client.Post(s.URL, "text/plain", strings.NewReader(tt.url))
		assert.NoError(t, err)
		defer postRes.Body.Close()
		body, err := io.ReadAll(postRes.Body)
		shortURL := string(body)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, postRes.StatusCode)
		assert.Contains(t, shortURL, serverAddress)
		// TODO: fix url construction
		shortURLPath := shortURL[len(serverAddress):]

		// do not follow redirect
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
		getRes, err := client.Get(s.URL + "/" + shortURLPath)
		require.NoError(t, err)
		defer getRes.Body.Close()
		_, err = io.ReadAll(getRes.Body)
		longURL := getRes.Header.Get("Location")
		assert.Equal(t, tt.url, longURL)
	}
}
