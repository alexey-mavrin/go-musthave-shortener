package app

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (c Config) storeJSONHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	var body []byte
	var err error

	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "NewReader "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()
		defer r.Body.Close()

		body, err = io.ReadAll(gz)
	} else {
		body, err = ioutil.ReadAll(r.Body)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
