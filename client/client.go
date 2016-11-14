package client

import (
	"net/http"
	"io/ioutil"
	"errors"
	"strconv"
	"fmt"
	"strings"
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
		return "", unexpectedStatusError(resp.StatusCode)
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)

	if (err != nil) {
		return "", err
	}

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
		return unexpectedStatusError(resp.StatusCode)
	}

	return nil

}

func update(url string) error  {

	req, err := http.NewRequest(http.MethodPatch, url, nil)

	if (err != nil) {
		return err
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if (err != nil) {
		return err
	}

	if (err != nil) {
		return err
	}

	if (resp.StatusCode == http.StatusNotFound){
		return ErrKeyNotFound
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return ErrServerError
	}

	if (resp.StatusCode != http.StatusOK) {
		return unexpectedStatusError(resp.StatusCode)
	}

	return nil

}

func (client *Client) Update(key string, value string) error {
	url := client.addr + "/keys?key=" + key + "&value=" + value
	return update(url)
}

func (client *Client) UpdateWithTtl(key string, value string, ttl int) error {
	url := client.addr + "/keys?key=" + key + "&value=" + value + "&ttl=" + strconv.Itoa(ttl)
	return update(url)
}


func (client *Client) Del(key string) error {
	url := client.addr + "/keys?key=" + key

	req, err := http.NewRequest(http.MethodDelete, url, nil)

	if (err != nil) {
		return err
	}

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)

	if (err != nil) {
		return err
	}

	if (resp.StatusCode == http.StatusNotFound){
		return ErrKeyNotFound
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return ErrServerError
	}

	if (resp.StatusCode != http.StatusOK) {
		return unexpectedStatusError(resp.StatusCode)
	}

	return nil
}

func unexpectedStatusError(status int) error {
	return fmt.Errorf("Unexpected status %d", status)
}

func (client *Client) Keys() ([]string, error) {

	url := client.addr + "/keys"

	resp, err := http.Get(url)

	if (err != nil) {
		return nil, err
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return  nil, ErrServerError
	}

	if (resp.StatusCode != http.StatusOK) {
		return nil, unexpectedStatusError(resp.StatusCode)
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)

	if (err != nil) {
		return nil, err
	}

	return strings.Split(string(content), "\n"), nil
}





