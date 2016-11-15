package handlers

import (
	"net/http"
	"gcache"
	"log"
)

type ListsHandler struct {
	Cache *gcache.Cache
}

func (handler *ListsHandler) Init(cache * gcache.Cache) Handler {
	return &ListsHandler{
		Cache: cache,
	}
}

func (handler *ListsHandler) Handle(w http.ResponseWriter, req *http.Request) {

	if err := req.ParseForm(); err != nil {
		log.Printf("Error parsing form: %s", err)
		return
	}



	// Nothing matched, return bad request
	w.WriteHeader(http.StatusBadRequest)
}