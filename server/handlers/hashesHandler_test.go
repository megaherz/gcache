package handlers

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"gcache"
	"io/ioutil"
)

func TestHashesHandler_HSetGSet(t *testing.T) {

	const hashKey  = "hashKey"
	const key  = "key"
	const value  = "value"

	handler := new(HashesHandler).Init(gcache.NewCache())

	ts := httptest.NewServer(http.HandlerFunc(handler.ServeHTTP))
	defer ts.Close()

	url := ts.URL + "?hashKey=" + hashKey + "&key=" + key + "&value=" + value

	rr, err := http.Post(url, "", nil)

	if (err != nil){
		t.Fatalf("http.Get(%q) unexpected error: %v", url, err)
	}

	// Check the status code is what we expect.
	if status := rr.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	url = ts.URL + "?hashKey=" + hashKey + "&key=" + key

	rr, err = http.Get(url)

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
	if value != string(actual) {
		t.Errorf("Expected the message '%s'\n", value)
	}

}
