package handlers

import (
	"net/http"
	"gcache"
	"encoding/csv"
	"bytes"
	"strings"
)

type Handler interface {
	Handle(w http.ResponseWriter, req *http.Request)
	Init(cache * gcache.Cache) Handler
}

func serialize(items []interface{}) (string, error) {

	var casted = make([]string, 0)

	for _, value := range items {
		casted = append(casted, value.(string))
	}

	return serializeStrings(casted)
}

func serializeStrings(items []string) (string, error) {

	buffer := &bytes.Buffer{} // creates IO Writer
	writer := csv.NewWriter(buffer)
	err := writer.Write(items)

	if err != nil {
		return "", err
	}

	writer.Flush()

	return strings.Trim(buffer.String(), "\n"), nil
}