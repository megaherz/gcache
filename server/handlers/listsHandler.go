package handlers

import (
	"net/http"
	"gcache"
	"log"
	"strings"
	"strconv"
	"fmt"
)

const (
	formListKey = "listKey"
	formRangeTo = "to"
	formRangeFrom = "from"
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

	path := strings.ToLower(req.URL.Path)

	// Command-Query Router
	switch req.Method {
	case http.MethodGet:
		if (strings.Contains(path, "range")) {
			handler.rangeQuery(w, req)
			return
		}

	case http.MethodPost:

		if (strings.Contains(path, "lpush")) {
			handler.lPushCommand(w, req)
			return
		} else if (strings.Contains(path, "rpush")) {
			handler.rPushCommand(w, req)
			return
		} else if (strings.Contains(path, "lpop")) {
			handler.lPopCommand(w, req)
			return
		} else if (strings.Contains(path, "rpop")) {
			handler.rPopCommand(w, req)
			return
		}
	}

	// Nothing matched, return bad request
	w.WriteHeader(http.StatusBadRequest)
}

func (handler *ListsHandler) rangeQuery(w http.ResponseWriter, req *http.Request) {

	listKey := req.Form.Get(formListKey)
	if (listKey == "") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	from, err := strconv.Atoi(req.Form.Get(formRangeFrom))
	if (err != nil) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	to, err := strconv.Atoi(req.Form.Get(formRangeTo))
	if (err != nil) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	items, err := handler.Cache.LRange(listKey, from, to)

	if (err != nil) {

		if (err == gcache.ErrKeyNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		fmt.Fprint(w, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Serialize items to string
	serialized, err := serialize(items)

	if (err != nil) {
		log.Printf("rangeQuery. Failed to serialize items %s", items)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, serialized)
}

func (handler *ListsHandler) lPushCommand(w http.ResponseWriter, req *http.Request) {

	listKey := req.Form.Get(formListKey)
	if (listKey == "") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value := req.Form.Get(formValue)
	if (value == "") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := handler.Cache.LPush(listKey, value)

	if (err != nil) {

		if (err == gcache.ErrKeyNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		fmt.Fprint(w, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (handler *ListsHandler) rPushCommand(w http.ResponseWriter, req *http.Request) {

	listKey := req.Form.Get(formListKey)
	if (listKey == "") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value := req.Form.Get(formValue)
	if (value == "") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := handler.Cache.RPush(listKey, value)

	if (err != nil) {

		if (err == gcache.ErrKeyNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		fmt.Fprint(w, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (handler *ListsHandler)  lPopCommand(w http.ResponseWriter, req *http.Request) {

	listKey := req.Form.Get(formListKey)
	if (listKey == "") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, err := handler.Cache.LPop(listKey)

	if (err != nil) {

		if (err == gcache.ErrKeyNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		fmt.Fprint(w, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, value)
}

func (handler *ListsHandler)  rPopCommand(w http.ResponseWriter, req *http.Request) {

	listKey := req.Form.Get(formListKey)
	if (listKey == "") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, err := handler.Cache.RPop(listKey)

	if (err != nil) {

		if (err == gcache.ErrKeyNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		fmt.Fprint(w, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, value)
}