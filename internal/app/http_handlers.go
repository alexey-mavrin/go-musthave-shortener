package app

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func optionalDecompressBody(r *http.Request) ([]byte, error) {
	if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return []byte{}, err
		}
		return body, nil
	}
	gz, err := gzip.NewReader(r.Body)
	if err != nil {
		return []byte{}, err
	}
	defer gz.Close()
	body, err := io.ReadAll(gz)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func (c Config) storeJSONHandler(w http.ResponseWriter, r *http.Request) {
	body, err := optionalDecompressBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	url := new(URL)
	err = json.Unmarshal(body, url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Print("url:", url.URL)
	key, err := c.sh.s.store(url.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := Result{
		Result: c.BaseURL + "/" + key,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(w)
	err = enc.Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c Config) storeHandler(w http.ResponseWriter, r *http.Request) {
	body, err := optionalDecompressBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	url := string(body)
	log.Print("url:", url)
	key, err := c.sh.s.store(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(c.BaseURL + "/" + key))
}

func (c Config) fetchHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	log.Print("key:", key)
	url, err := c.sh.s.get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
