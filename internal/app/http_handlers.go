package app

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const serverAddress = `http://localhost:8080/`

func (sh store) storeHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	url := string(body)
	log.Print("url:", url)
	key, err := sh.s.store(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(serverAddress + key))
}

func (sh store) fetchHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	log.Print("key:", key)
	url, err := sh.s.get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
