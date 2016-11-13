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

	if (key == "") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch req.Method {
	case http.MethodPost:

		// Parse value
		value := req.Form.Get("value")

		if (value == "") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		strTtl := req.Form.Get("ttl")

		if (strTtl == "") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ttl, err := strconv.Atoi(strTtl)

		if (err != nil) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		handler.set(key, value, ttl)

	case http.MethodGet:

		handler.get(key, w)

	case http.MethodPatch:

		value := req.Form.Get("value")

		if (value == "") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Parse ttl
		if (req.Form.Get("ttl") == "") {
			handler.update(key, value, noTtlDefined, w)
		} else {
			ttl, err := strconv.Atoi(req.Form.Get("ttl"))
			if err != nil || ttl < 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			handler.update(key, value, ttl, w)
		}



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
	handler.cache.Set(key, value, handler.intToDurationInMinutes(ttl))
}

func (handler *keysHandler) intToDurationInMinutes(ttl int) time.Duration {
	return time.Duration(float64(int64(ttl) * time.Second.Nanoseconds()))
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
		err = handler.cache.UpdateWithTll(key, value, handler.intToDurationInMinutes(ttl))
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
