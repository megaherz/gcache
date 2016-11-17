package client

import (
	"net/http"
	"io/ioutil"
	"errors"
	"strconv"
	"fmt"
	"encoding/csv"
	"io"
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

func NewClientWithAuth (addr string, psw string) *Client  {
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

	resp, err := http.Get(fmt.Sprintf("%s/keys", client.addr))

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

func (client*Client) LPush(listKey string, value string) error {
	return client.push("lpush", listKey, value)
}

func (client*Client) RPush(listKey string, value string) error {
	return client.push("rpush", listKey, value)
}

func (client*Client) push(method string, listKey string, value string) error {
	resp, err := http.Post(fmt.Sprintf("%s/%s?listKey=%s&value=%s", client.addr, method, listKey, value), "", nil)

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

func (client* Client) pop(method string, listKey string) (string, error) {
	resp, err := http.Post(fmt.Sprintf("%s/%s?listKey=%s", client.addr, method, listKey), "", nil)

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

func (client*Client) LPop(listKey string) (string, error) {
	return client.pop("lpop", listKey)
}

func (client*Client) RPop(listKey string) (string, error) {
	return client.pop("rpop", listKey)
}

func (client*Client) LRange(listKey string, from int, to int) ([]string, error) {

	url := fmt.Sprintf("%s/range?listKey=%s&from=%d&to=%d", client.addr, listKey, from, to)

	resp, err := http.Get(url)

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

func readCsv(body io.Reader) ([]string, error) {
	reader := csv.NewReader(body)
	result, err := reader.Read()

	if (err != nil && err != io.EOF) {
		return nil, err
	}

	return result, nil
}

//------- HASH -----------

func (client* Client) HGet(hashKey string, key string) (string, error) {
	url := fmt.Sprintf("%s/hashes?hashKey=%s&key=%s", client.addr, hashKey, key)

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

func (client* Client) HSet(hashKey string, key string, value string) error  {
	url := fmt.Sprintf("%s/hashes?hashKey=%s&key=%s&value=%s", client.addr, hashKey, key, value)

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







