package handlers

import (
	"fmt"
	"gcache"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	formTtl = "ttl"
)

type KeysHandler struct {
	Cache *gcache.Cache
}

func (handler *KeysHandler) Init(cache *gcache.Cache) Handler {
	return &KeysHandler{
		Cache: cache,
	}
}

func (handler *KeysHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	if err := req.ParseForm(); err != nil {
		log.Printf("Error parsing form: %s", err)
		return
	}

	// Keys
	if len(req.Form) == 0 && req.Method == http.MethodGet {
		handler.keysQuery(w, req)
		return
	}

	switch req.Method {

	// Get
	case http.MethodGet:
		handler.getKeyQuery(w, req)
		return

	// Set
	case http.MethodPost:
		handler.setKeyCommand(w, req)
		return

	// Update
	case http.MethodPatch:
		handler.updateCommand(w, req)
		return

	// Delete
	case http.MethodDelete:
		handler.removeCommand(w, req)
		return
	}

	// Nothing matched, return bad request
	w.WriteHeader(http.StatusBadRequest)
}

func (handler *KeysHandler) keysQuery(w http.ResponseWriter, req *http.Request) {

	keys := handler.Cache.Keys()

	// Serialize keys to string
	serialized, err := serializeStrings(keys)

	if err != nil {
		log.Printf("keysQuery. Failed to serialize keys %s", keys)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, serialized)
}

func (handler *KeysHandler) getKeyQuery(w http.ResponseWriter, req *http.Request) {

	key := req.Form.Get(formKey)

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, err := handler.Cache.Get(key)

	if err == gcache.ErrKeyNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, value.(string))
}

func (handler *KeysHandler) setKeyCommand(w http.ResponseWriter, req *http.Request) {

	key := req.Form.Get(formKey)

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value := req.Form.Get(formValue)

	if value == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	strTtl := req.Form.Get(formTtl)

	if strTtl == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ttl, err := strconv.Atoi(strTtl)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	handler.Cache.Set(key, value, convertIntToDurationInMinutes(ttl))
}

func (handler *KeysHandler) removeCommand(w http.ResponseWriter, req *http.Request) {

	key := req.Form.Get(formKey)

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := handler.Cache.Del(key)

	if err == gcache.ErrKeyNotFound {
		http.NotFound(w, req)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *KeysHandler) updateCommand(w http.ResponseWriter, req *http.Request) {

	key := req.Form.Get(formKey)

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value := req.Form.Get(formValue)

	if value == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}


	if req.Form.Get(formTtl) == "" {
		err := handler.Cache.Update(key, value)

		if err == gcache.ErrKeyNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		ttl, err := strconv.Atoi(req.Form.Get(formTtl))
		if err != nil || ttl < 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = handler.Cache.UpdateWithTll(key, value, convertIntToDurationInMinutes(ttl))

		if err == gcache.ErrKeyNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

}

// Convert int to duration in minutes
func convertIntToDurationInMinutes(ttl int) time.Duration {
	return time.Duration(float64(int64(ttl) * time.Second.Nanoseconds()))
}
