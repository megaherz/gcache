package handlers

import (
	"gcache"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestKeysHandler_SetGet(t *testing.T) {

	const key = "key1"
	const value = "value"

	keysHandler := KeysHandler{
		Cache: gcache.NewCache(),
	}

	ts := httptest.NewServer(http.HandlerFunc(keysHandler.ServeHTTP))
	defer ts.Close()

	url := ts.URL + "?key=" + key

	rr, err := http.Get(url)

	if err != nil {
		t.Fatalf("http.Get(%q) unexpected error: %v", url, err)
	}

	// Check the status code is what we expect.
	if status := rr.StatusCode; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	url = ts.URL + "?key=" + key + "&ttl=5&value=" + value
	rr, err = http.Post(url, "", nil)

	if err != nil {
		t.Fatalf("http.Post(%q) unexpected error: %v", url, err)
	}

	// Check the status code is what we expect.
	if status := rr.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Bad request. Ttl is missed
	url = ts.URL + "?key=" + key + "&value=" + value
	rr, err = http.Post(url, "", nil)

	if err != nil {
		t.Fatalf("http.Post(%q) unexpected error: %v", url, err)
	}

	// Check the status code is what we expect.
	if status := rr.StatusCode; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	rr, err = http.Get(url)

	if err != nil {
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

func TestKeysHandler_Keys(t *testing.T) {
	keysHandler := KeysHandler{
		Cache: gcache.NewCache(),
	}

	ts := httptest.NewServer(http.HandlerFunc(keysHandler.ServeHTTP))
	defer ts.Close()

	url := ts.URL

	rr, err := http.Get(url)

	if err != nil {
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
	if "" != string(actual) {
		t.Errorf("Expected no keys but recieved '%s'", string(actual))
	}

	// Insert some keys
	const n = 10
	for i := 0; i < n; i++ {
		key := strconv.Itoa(i)
		http.Post(ts.URL+"?key="+key+"&ttl=5&value=value", "", nil)
	}

	rr, err = http.Get(url)

	if err != nil {
		t.Fatalf("http.Get(%q) unexpected error: %v", url, err)
	}

	// Check the status code is what we expect.
	if status := rr.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	actual, err = ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("Actual", string(actual))
}

func TestKeysHandler_SetDel(t *testing.T) {
	const key = "key1"
	const value = "value"

	keysHandler := KeysHandler{
		Cache: gcache.NewCache(),
	}

	ts := httptest.NewServer(http.HandlerFunc(keysHandler.ServeHTTP))
	defer ts.Close()

	url := ts.URL + "?key=" + key + "&ttl=5&value=" + value
	rr, err := http.Post(url, "", nil)

	if err != nil {
		t.Fatalf("http.Post(%q) unexpected error: %v", url, err)
	}

	// Check the status code is what we expect.
	if status := rr.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	url = ts.URL + "?key=" + key
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	rr, err = http.DefaultClient.Do(req)

	if err != nil {
		t.Fatalf("http.Post(%q) unexpected error: %v", url, err)
	}

	// Check the status code is what we expect.
	if status := rr.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}
