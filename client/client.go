package client

import (
	"net/http"
	"io/ioutil"
	"errors"
	"fmt"
	"encoding/csv"
	"io"
)

const headerAuthorization  = "Authorization"

var ErrKeyNotFound = errors.New("Key Not Found")
var ErrServerError = errors.New("Internal Server error")

type Client struct {
	addr string
	psw string
}

func NewClient (addr string) *Client  {
	return NewClientWithAuth(addr, "")
}

func NewClientWithAuth (addr string, psw string) *Client  {
	return &Client{
		addr: addr,
		psw: psw,
	}
}

func (client *Client) doRequest(method, urlStr string, body io.Reader) (*http.Response, error) {

	req, err := http.NewRequest(method, client.addr + urlStr, body)
	if (err != nil) {
		return nil, err
	}

	// Set the authorization header if password is set
	if (client.psw != "") {
		req.Header.Set(headerAuthorization, client.psw)
	}

	return http.DefaultClient.Do(req)
}


func (client *Client) Get(key string) (string, error) {

   	resp, err := client.doRequest(http.MethodGet, "/keys?key=" + key, nil)

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

	url := fmt.Sprintf("/keys?key=%s&value=%s&ttl=%d", key, value, ttl)

	resp, err := client.doRequest(http.MethodPost, url, nil)
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

func (client *Client) updateKey(url string) error  {

	resp, err := client.doRequest(http.MethodPatch, url, nil)

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
	url := fmt.Sprintf("/keys?key=%s&value=%s", key, value)
	return client.updateKey(url)
}

func (client *Client) UpdateWithTtl(key string, value string, ttl int) error {
	url := fmt.Sprintf("/keys?key=%s&value=%s&ttl=%d", key, value, ttl)
	return client.updateKey(url)
}


func (client *Client) Del(key string) error {
	url := "/keys?key=" + key

	resp, err := client.doRequest(http.MethodDelete, url, nil)

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

	resp, err := client.doRequest(http.MethodGet, "/keys", nil)

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

	return readCsv(resp.Body)
}

//------ LIST ---------

func (client*Client) LPush(key string, value string) error {
	return client.push("lpush", key, value)
}

func (client*Client) RPush(key string, value string) error {
	return client.push("rpush", key, value)
}


func (client*Client) LRange(key string, from int, to int) ([]string, error) {

	url := fmt.Sprintf("/lists?op=range&key=%s&from=%d&to=%d", key, from, to)

	resp, err := client.doRequest(http.MethodGet, url, nil)

	if (err != nil) {
		return nil, err
	}

	if (resp.StatusCode == http.StatusNotFound){
		return nil, ErrKeyNotFound
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return nil, ErrServerError
	}

	if (resp.StatusCode != http.StatusOK) {
		return nil, unexpectedStatusError(resp.StatusCode)
	}

	defer resp.Body.Close()

	return readCsv(resp.Body)
}

func (client*Client) LPop(key string) (string, error) {
	return client.pop("lpop", key)
}

func (client*Client) RPop(key string) (string, error) {
	return client.pop("rpop", key)
}

func (client*Client) push(method string, key string, value string) error {

	url := fmt.Sprintf("/lists?op=%s&key=%s&value=%s",method, key, value)

	resp, err := client.doRequest(http.MethodPost, url, nil)

	if (err != nil){
		return err
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return  ErrServerError
	}

	if (resp.StatusCode != http.StatusOK) {
		return unexpectedStatusError(resp.StatusCode)
	}

	return nil
}

func (client* Client) pop(method string, key string) (string, error) {

	url := fmt.Sprintf("/lists?op=%s&key=%s", method, key)
	resp, err := client.doRequest(http.MethodPost, url, nil)

	if (err != nil){
		return "", err
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return  "", ErrServerError
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


func readCsv(body io.Reader) ([]string, error) {
	reader := csv.NewReader(body)
	result, err := reader.Read()

	if (err != nil && err != io.EOF) {
		return nil, err
	}

	return result, nil
}

//------- HASH -----------

func (client* Client) HGet(key string, hashKey string) (string, error) {
	url := fmt.Sprintf("/hashes?key=%s&hashKey=%s", key, hashKey)

	resp, err := client.doRequest(http.MethodGet, url, nil)

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

func (client* Client) HSet(key string, hashKey string, value string) error  {
	url := fmt.Sprintf("/hashes?key=%s&hashKey=%s&value=%s", key, hashKey, value)

	resp, err := client.doRequest(http.MethodPost, url, nil)

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







