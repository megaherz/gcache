package handlers

import (
	"net/http"
	"gcache"
)

type ListHandler struct {
	Cache *gcache.Cache
}


func (handler *ListHandler) Handle(w http.ResponseWriter, req *http.Request) {
}