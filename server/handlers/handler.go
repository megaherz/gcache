package handlers

import (
	"net/http"
	"gcache"
	"strings"
)

type Handler interface {
	Handle(w http.ResponseWriter, req *http.Request)
	Init(cache * gcache.Cache) Handler
}

func serialize(items []interface{}) string {

	var casted = make([]string, 0)

	for _, value := range items {
		casted = append(casted, value.(string))
	}

	return serializeStrings(casted)
}

func serializeStrings(items []string) string {
	return strings.Join(items, "\n")
}