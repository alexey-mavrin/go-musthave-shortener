package app

import (
	"io/ioutil"
	"log"
	"net/http"
)

const serverAddress = `http://localhost:8080/`

func (sh storeHandler) storeHandler(w http.ResponseWriter, r *http.Request) {
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

func (sh storeHandler) fetchHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.EscapedPath()[1:]
	log.Print("key:", key)
	url, err := sh.s.get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
