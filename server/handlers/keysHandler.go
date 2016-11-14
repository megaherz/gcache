package handlers

import (
	"net/http"
	"log"
	"gcache"
	"fmt"
	"time"
	"strconv"
	"strings"
)

type KeysHandler struct {
   Cache *gcache.Cache
}

const noTtlDefined int = -1

func (handler *KeysHandler) Init(cache * gcache.Cache) Handler {
	return &KeysHandler{
		Cache: cache,
	}
}

func (handler *KeysHandler) Handle(w http.ResponseWriter, req *http.Request) {

	if err := req.ParseForm(); err != nil {
		log.Printf("Error parsing form: %s", err)
		return
	}

	// Keys
	if (len(req.Form) == 0 && req.Method == http.MethodGet) {
		handler.keys(w)
		return;
	}

	key := req.Form.Get("key")

	if (key == "") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch req.Method {

	// Get
	case http.MethodGet:
		handler.get(key, w)
		return

	// Set
	case http.MethodPost:

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
		return

	// Update
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

		return


	// Delete
	case http.MethodDelete:
		handler.remove(key, w, req)
		return
	}

	// Nothing matched, return bad request
	w.WriteHeader(http.StatusBadRequest)
}

func (handler *KeysHandler) keys(w http.ResponseWriter) {
	keys := handler.Cache.Keys()
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, strings.Join(keys, " "))
}

func (handler *KeysHandler) get(key string, w http.ResponseWriter){

	value, err := handler.Cache.Get(key)

	if (err == gcache.ErrKeyNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, value.(string))
}

func (handler *KeysHandler) set(key string, value string, ttl int){
	handler.Cache.Set(key, value, handler.intToDurationInMinutes(ttl))
}

func (handler *KeysHandler) intToDurationInMinutes(ttl int) time.Duration {
	return time.Duration(float64(int64(ttl) * time.Second.Nanoseconds()))
}
func (handler *KeysHandler) remove(key string, w http.ResponseWriter, req *http.Request) {
	err := handler.Cache.Del(key)

	if (err == gcache.ErrKeyNotFound) {
		http.NotFound(w, req)
		return
	}

	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *KeysHandler) update(key string, value string, ttl int, w http.ResponseWriter) {

	var err error

	if (ttl != noTtlDefined) {
		err = handler.Cache.Update(key, value)
	} else {
		err = handler.Cache.UpdateWithTll(key, value, handler.intToDurationInMinutes(ttl))
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
