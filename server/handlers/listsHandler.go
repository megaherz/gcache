package handlers

import (
	"fmt"
	"gcache"
	"log"
	"net/http"
	"strconv"
)

const (
	formRangeTo   = "to"
	formRangeFrom = "from"
	formOperation = "op"
)

type ListsHandler struct {
	Cache *gcache.Cache
}

func (handler *ListsHandler) Init(cache *gcache.Cache) Handler {
	return &ListsHandler{
		Cache: cache,
	}
}

func (handler *ListsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	if err := req.ParseForm(); err != nil {
		log.Printf("Error parsing form: %s", err)
		return
	}

	operation := req.Form.Get(formOperation)

	// Command-Query Router
	switch req.Method {
	case http.MethodGet:
		if operation == "range" {
			handler.rangeQuery(w, req)
			return
		}

	case http.MethodPost:

		if operation == "lpush" {
			handler.lPushCommand(w, req)
			return
		} else if operation == "rpush" {
			handler.rPushCommand(w, req)
			return
		} else if operation == "lpop" {
			handler.lPopCommand(w, req)
			return
		} else if operation == "rpop" {
			handler.rPopCommand(w, req)
			return
		}
	}

	// Nothing matched, return bad request
	w.WriteHeader(http.StatusBadRequest)
}

func (handler *ListsHandler) rangeQuery(w http.ResponseWriter, req *http.Request) {

	key := req.Form.Get(formKey)
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	from, err := strconv.Atoi(req.Form.Get(formRangeFrom))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	to, err := strconv.Atoi(req.Form.Get(formRangeTo))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	items, err := handler.Cache.LRange(key, from, to)

	if err != nil {

		if err == gcache.ErrKeyNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		fmt.Fprint(w, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Serialize items to string
	serialized, err := serialize(items)

	if err != nil {
		log.Printf("rangeQuery. Failed to serialize items %s", items)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, serialized)
}

func (handler *ListsHandler) lPushCommand(w http.ResponseWriter, req *http.Request) {

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

	err := handler.Cache.LPush(key, value)

	if err != nil {
		fmt.Fprint(w, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (handler *ListsHandler) rPushCommand(w http.ResponseWriter, req *http.Request) {

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

	err := handler.Cache.RPush(key, value)

	if err != nil {
		fmt.Fprint(w, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (handler *ListsHandler) lPopCommand(w http.ResponseWriter, req *http.Request) {

	key := req.Form.Get(formKey)
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, err := handler.Cache.LPop(key)

	if err != nil {

		if err == gcache.ErrKeyNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		fmt.Fprint(w, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, value)
}

func (handler *ListsHandler) rPopCommand(w http.ResponseWriter, req *http.Request) {

	key := req.Form.Get(formKey)
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, err := handler.Cache.RPop(key)

	if err != nil {

		if err == gcache.ErrKeyNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		fmt.Fprint(w, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, value)
}
