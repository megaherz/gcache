package handlers

import (
	"testing"
	"gcache"
	"net/http/httptest"
	"net/http"
)

func TestListsHandler_LPush_LPop(t *testing.T) {

	const listKey  = "listKey"
	const value  = "value"

	handler := new(ListsHandler).Init(gcache.NewCache())

	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	url := ts.URL + "/lpush?listKey=" + listKey +  "&value=" + value

	rr, err := http.Post(url, "", nil)

	if (err != nil){
		t.Fatalf("http.Post(%q) unexpected error: %v", url, err)
	}

	// Check the status code is what we expect.
	if status := rr.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
