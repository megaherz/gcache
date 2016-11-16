package handlers

import (
	"testing"
	"gcache"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"fmt"
	"log"
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

	url = ts.URL + "/lpop?listKey=" + listKey

	rr, err = http.Post(url, "", nil)

	if (err != nil){
		t.Fatalf("http.Post(%q) unexpected error: %v", url, err)
	}

	// Check the status code is what we expect.
	if status := rr.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	actual, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}
	if value != string(actual) {
		t.Errorf("Expected the message '%s'\n", value)
	}

}

func TestListsHandler_Range(t *testing.T) {

	const listKey  = "listKey"

	handler := new(ListsHandler).Init(gcache.NewCache())

	ts := httptest.NewServer(http.HandlerFunc(handler.Handle))
	defer ts.Close()

	// LPUSH 10 items
	for i := 0; i < 10; i++ {
		url := fmt.Sprintf("%s/lpush?listKey=%s&value=%d", ts.URL, listKey, i)
		_, err := http.Post(url, "", nil)
		if (err != nil) {
			t.Errorf("Failed to lpush %url", url)
		}
	}

	url := fmt.Sprintf("%s/range?listKey=%s&from=%d&to=%d", ts.URL, listKey, 2, 4)

	rr, err := http.Get(url)

	if (err != nil){
		t.Fatalf("http.Get(%q) unexpected error: %v", url, err)
	}

	// Check the status code is what we expect.
	if status := rr.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	actual, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("RANGE response", string(actual))

	value := "2,3,4"

	if value != string(actual) {
		t.Errorf("Expected the message '%s' but recieved '%s'\n", value, string(actual))
	}

}
