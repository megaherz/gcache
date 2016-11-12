package server

import (
	"net/http"
	"log"
	"gcache"
	"fmt"
	"time"
	"strconv"
)

type keysHandler struct {
   cache *gcache.Cache
}

const noTtlDefined int = -1

func (handler *keysHandler) Handle(w http.ResponseWriter, req *http.Request) {

	if err := req.ParseForm(); err != nil {
		log.Printf("Error parsing form: %s", err)
		return
	}

	key := req.Form.Get("key")

	if (key == nil) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch req.Method {
	case http.MethodPost:

		// Parse value
		value := req.Form["value"]

		if (value == nil) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Parse ttl
		var ttl int
		if (req.Form["ttl"] == nil) {
			ttl = noTtlDefined
		} else {
			ttl, err := strconv.Atoi(req.Form["ttl"])
			if err != nil || ttl < 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		handler.set(key, value, strconv.Atoi(ttl))

	case http.MethodGet:

		handler.get(key, w)

	case http.MethodPatch:

		value := req.Form["value"]

		if (value == nil) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Parse ttl
		var ttl int
		if (req.Form["ttl"] == nil) {
			ttl = noTtlDefined
		} else {
			ttl, err := strconv.Atoi(req.Form["ttl"])
			if err != nil || ttl < 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		handler.update(key, value,  strconv.Atoi(ttl), w)

	case http.MethodDelete:

		handler.remove(key, w, req)
	}


}

func (handler *keysHandler) get(key string, w http.ResponseWriter){
	value, err := handler.cache.Get(key)

	if (err == gcache.ErrKeyNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, value.(string))
}

func (handler *keysHandler) set(key string, value string, ttl int){
	handler.cache.Set(key, value, ttl * time.Second)
}

func (handler *keysHandler) remove(key string, w http.ResponseWriter, req *http.Request) {
	err := handler.cache.Del(key)

	if (err == gcache.ErrKeyNotFound) {
		http.NotFound(w, req)
		return
	}

	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *keysHandler) update(key string, value string, ttl int, w http.ResponseWriter) {

	var err error

	if (ttl != noTtlDefined) {
		err = handler.cache.Update(key, value)
	} else {
		err = handler.cache.UpdateWithTll(key, value, ttl)
	}

	if (err == gcache.ErrKeyNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
