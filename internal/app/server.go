package app

import (
	"net/http"
)

func (sh storeHandler) dispatchHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		sh.fetchHandler(w, r)
	case http.MethodPost:
		sh.storeHandler(w, r)
	}
}

func Run() error {
	sh := newStoreHandler()
	http.HandleFunc("/", sh.dispatchHandler)
	return http.ListenAndServe(":8080", nil)
}
