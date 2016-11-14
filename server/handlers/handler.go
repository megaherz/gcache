package handlers

import (
	"net/http"
	"gcache"
)

type Handler interface {
	Handle(w http.ResponseWriter, req *http.Request)
	Init(cache * gcache.Cache) Handler
}
