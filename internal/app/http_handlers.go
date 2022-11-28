package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const serverAddress = `http://localhost:8080/`

func (sh store) storeJSONHandler(w http.ResponseWriter, r *http.Request) {
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
	key, err := sh.s.store(url.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	res := Result{
		Result: serverAddress + key,
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

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
