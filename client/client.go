package client

import (
	"net/http"
	"io/ioutil"
	"errors"
	"strconv"
)

var ErrKeyNotFound = errors.New("Key Not Found")
var ErrServerError = errors.New("Internal Server error")

type Client struct {
	addr string
}

func NewClient (addr string) *Client  {
	return &Client{
		addr: addr,
	}
}

func (client *Client) Get(key string) (string, error) {

	url := client.addr + "/keys?key=" + key

   	resp, err := http.Get(url)

	if (err != nil) {
		return "", err
	}

	if (resp.StatusCode == http.StatusNotFound){
		return "", ErrKeyNotFound
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return "", ErrServerError
	}

	if (resp.StatusCode != http.StatusOK) {
		return "", errors.New("Unexpected status " + strconv.Itoa(resp.StatusCode))
	}

	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return string(content), nil

}

func (client *Client) Set(key string, value string, ttl int) error {

	url := client.addr + "/keys?key=" + key + "&value=" + value + "&ttl=" + strconv.Itoa(ttl)

	resp, err := http.Post(url, "", nil)
	if (err != nil) {
		return err
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return ErrServerError
	}

	if (resp.StatusCode != http.StatusOK) {
		return errors.New("Unexpected status " + strconv.Itoa(resp.StatusCode))
	}

	return nil

}



