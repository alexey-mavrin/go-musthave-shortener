package app

import (
	"log"
	"net/http"
)

const serverAddress = `http://localhost:8080/`

func (s *storage) storeHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	log.Print("url:", url)
	key, err := s.store(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(serverAddress + key))
}

func (s *storage) fetchHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.EscapedPath()[1:]
	log.Print("key:", key)
	url, err := s.get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *storage) dispatchHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.fetchHandler(w, r)
	case http.MethodPost:
		s.storeHandler(w, r)
	}
}

func Run() error {
	storage := newStorage()
	http.HandleFunc("/", storage.dispatchHandler)
	return http.ListenAndServe(":8080", nil)
}
