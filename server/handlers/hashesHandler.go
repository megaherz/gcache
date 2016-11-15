package handlers

import (
	"net/http"
	"gcache"
	"log"
	"fmt"
)

type HashesHandler struct {
	Cache *gcache.Cache
}

func (handler *HashesHandler) Init(cache * gcache.Cache) Handler {
	return &HashesHandler{
		Cache: cache,
	}
}

func (handler *HashesHandler) Handle(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		log.Printf("Error parsing form: %s", err)
		return
	}

	key := req.Form.Get("key")

	if (key == "") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashKey := req.Form.Get("hashKey")

	if (hashKey == "") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}


	switch req.Method {
	case http.MethodGet:
		value, err := handler.Cache.HGet(hashKey, key)
		if (err != nil) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, value.(string))
		return

	case http.MethodPost:

		value := req.Form.Get("value")

		if (value == "") {
			w.WriteHeader(http.StatusBadRequest)

		}

		handler.Cache.HSet(hashKey, key, value)
		return
	}

	// Nothing matched, return bad request
	w.WriteHeader(http.StatusBadRequest)
}