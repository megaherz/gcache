package handlers

import (
	"bytes"
	"encoding/csv"
	"gcache"
	"net/http"
	"strings"
)

type Handler interface {
	http.Handler
	Init(cache *gcache.Cache) Handler
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
