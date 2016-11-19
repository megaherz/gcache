package handlers

import (
	"fmt"
	"gcache"
	"log"
	"net/http"
)

const (
	formHashKey = "hashKey"
	formKey     = "key"
	formValue   = "value"
)

type HashesHandler struct {
	Cache *gcache.Cache
}

func (handler *HashesHandler) Init(cache *gcache.Cache) Handler {
	return &HashesHandler{
		Cache: cache,
	}
}

func (handler *HashesHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		log.Printf("Error parsing form: %s", err)
		return
	}

	// Router
	switch req.Method {
	case http.MethodGet:
		handler.hGetQuery(w, req)
		return

	case http.MethodPost:
		handler.hSetCommand(w, req)
		return
	}

	// Nothing matched, return Bad request
	w.WriteHeader(http.StatusBadRequest)
}

func (handler *HashesHandler) hSetCommand(w http.ResponseWriter, req *http.Request) {
	hashKey := req.Form.Get(formHashKey)

	if hashKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	key := req.Form.Get(formKey)

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value := req.Form.Get(formValue)

	if value == "" {
		w.WriteHeader(http.StatusBadRequest)

	}

	err := handler.Cache.HSet(key, hashKey, value)

	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (handler *HashesHandler) hGetQuery(w http.ResponseWriter, req *http.Request) {
	hashKey := req.Form.Get(formHashKey)

	if hashKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	key := req.Form.Get(formKey)

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, err := handler.Cache.HGet(key, hashKey)

	w.Header().Set("Content-Type", "text/plain")

	if err != nil {

		if err == gcache.ErrKeyNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Key not found")
			return
		}

		if err == gcache.ErrHashKeyNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Hash key not found")
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	fmt.Fprint(w, value.(string))

}
