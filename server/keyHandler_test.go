package server

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"gcache"
	"io/ioutil"
)

func TestKeysHandler_SetGet(t *testing.T) {

	const key = "key1"
	const value  = "value"

	keysHandler := keysHandler{
		cache: gcache.NewCache(),
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	ts := httptest.NewServer(http.HandlerFunc(keysHandler.Handle))
	defer ts.Close()

	url := ts.URL + "?key=" + key

	rr, err := http.Get(url)

	if (err != nil){
		t.Fatalf("http.Get(%q) unexpected error: %v", url, err)
	}

	// Check the status code is what we expect.
	if status := rr.StatusCode; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	rr, err = http.Post(ts.URL + "?key=" + key + "&ttl=5&value=" + value, "", nil)

	if (err != nil){
		t.Fatalf("http.Post(%q) unexpected error: %v", url, err)
	}

	// Check the status code is what we expect.
	if status := rr.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Bad request. Ttl is missed
	rr, err = http.Post(ts.URL + "?key=" + key + "&value=" + value, "", nil)

	if (err != nil){
		t.Fatalf("http.Post(%q) unexpected error: %v", url, err)
	}

	// Check the status code is what we expect.
	if status := rr.StatusCode; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

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
