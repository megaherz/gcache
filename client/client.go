package client

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const headerAuthorization = "Authorization"

var ErrKeyNotFound = errors.New("Key Not Found")
var ErrServerError = errors.New("Internal Server error")

type Client struct {
	conns Connections
}

func NewClient(conns Connections) *Client {
	return &Client{
		conns: conns,
	}
}

func (client *Client) Get(key string) (string, error) {

	conn := client.conns.getShard(key)
	resp, err := conn.doRequest(http.MethodGet, "/keys?key=" + key, nil)

	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusNotFound {
		return "", ErrKeyNotFound
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return "", ErrServerError
	}

	if resp.StatusCode != http.StatusOK {
		return "", unexpectedStatusError(resp.StatusCode)
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(content), nil

}

func (client *Client) Set(key string, value string, ttl int) error {

	url := fmt.Sprintf("/keys?key=%s&value=%s&ttl=%d", key, value, ttl)

	conn := client.conns.getShard(key)
	resp, err := conn.doRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return ErrServerError
	}

	if resp.StatusCode != http.StatusOK {
		return unexpectedStatusError(resp.StatusCode)
	}

	return nil

}

func (client *Client) updateKey(conn Connection, url string) error {

	resp, err := conn.doRequest(http.MethodPatch, url, nil)

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusNotFound {
		return ErrKeyNotFound
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return ErrServerError
	}

	if resp.StatusCode != http.StatusOK {
		return unexpectedStatusError(resp.StatusCode)
	}

	return nil

}

func (client *Client) Update(key string, value string) error {
	conn := client.conns.getShard(key)
	url := fmt.Sprintf("/keys?key=%s&value=%s", key, value)
	return client.updateKey(conn, url)
}

func (client *Client) UpdateWithTtl(key string, value string, ttl int) error {
	conn := client.conns.getShard(key)
	url := fmt.Sprintf("/keys?key=%s&value=%s&ttl=%d", key, value, ttl)
	return client.updateKey(conn, url)
}

func (client *Client) Del(key string) error {

	conn := client.conns.getShard(key)

	url := "/keys?key=" + key

	resp, err := conn.doRequest(http.MethodDelete, url, nil)

	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusNotFound {
		return ErrKeyNotFound
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return ErrServerError
	}

	if resp.StatusCode != http.StatusOK {
		return unexpectedStatusError(resp.StatusCode)
	}

	return nil
}

func (client *Client) Keys() ([]string, error) {

	responses := client.conns.doParallelGetRequest("/keys")

	// If there are any error, return it
	for _, resp := range responses {
		if resp.err != nil {
			return nil, resp.err
		}
	}

	// Concatenate keys
	keys := []string{}
	for _, resp := range responses {
		content, err := readCsv(resp.response.Body)

		if err != nil {
			return nil, err
		}
		keys = append(keys, content...)
	}

	return keys, nil
}

//------ LIST ---------

func (client *Client) LPush(key string, value string) error {
	return client.push("lpush", key, value)
}

func (client *Client) RPush(key string, value string) error {
	return client.push("rpush", key, value)
}

func (client *Client) LRange(key string, from int, to int) ([]string, error) {


	url := fmt.Sprintf("/lists?op=range&key=%s&from=%d&to=%d", key, from, to)

	conn := client.conns.getShard(key)
	resp, err := conn.doRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrKeyNotFound
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return nil, ErrServerError
	}

	if resp.StatusCode != http.StatusOK {
		return nil, unexpectedStatusError(resp.StatusCode)
	}

	defer resp.Body.Close()

	return readCsv(resp.Body)
}

func (client *Client) LPop(key string) (string, error) {
	return client.pop("lpop", key)
}

func (client *Client) RPop(key string) (string, error) {
	return client.pop("rpop", key)
}

func (client *Client) push(method string, key string, value string) error {

	url := fmt.Sprintf("/lists?op=%s&key=%s&value=%s", method, key, value)

	conn := client.conns.getShard(key)
	resp, err := conn.doRequest(http.MethodPost, url, nil)

	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return ErrServerError
	}

	if resp.StatusCode != http.StatusOK {
		return unexpectedStatusError(resp.StatusCode)
	}

	return nil
}

func (client *Client) pop(method string, key string) (string, error) {

	url := fmt.Sprintf("/lists?op=%s&key=%s", method, key)

	conn := client.conns.getShard(key)
	resp, err := conn.doRequest(http.MethodPost, url, nil)

	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return "", ErrServerError
	}

	if resp.StatusCode != http.StatusOK {
		return "", unexpectedStatusError(resp.StatusCode)
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(content), nil
}

func unexpectedStatusError(status int) error {
	return fmt.Errorf("Unexpected status %d", status)
}

func readCsv(body io.Reader) ([]string, error) {
	reader := csv.NewReader(body)
	result, err := reader.Read()

	if err != nil && err != io.EOF {
		return nil, err
	}

	return result, nil
}

//------- HASH -----------

func (client *Client) HGet(key string, hashKey string) (string, error) {
	url := fmt.Sprintf("/hashes?key=%s&hashKey=%s", key, hashKey)

	conn := client.conns.getShard(key)
	resp, err := conn.doRequest(http.MethodGet, url, nil)

	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusNotFound {
		return "", ErrKeyNotFound
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return "", ErrServerError
	}

	if resp.StatusCode != http.StatusOK {
		return "", unexpectedStatusError(resp.StatusCode)
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (client *Client) HSet(key string, hashKey string, value string) error {
	url := fmt.Sprintf("/hashes?key=%s&hashKey=%s&value=%s", key, hashKey, value)

	conn := client.conns.getShard(key)
	resp, err := conn.doRequest(http.MethodPost, url, nil)

	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return ErrServerError
	}

	if resp.StatusCode != http.StatusOK {
		return unexpectedStatusError(resp.StatusCode)
	}

	return nil
}
