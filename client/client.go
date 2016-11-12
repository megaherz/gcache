package client

import (
	"net/http"
	"io/ioutil"
)

type Client struct {
	addr string
}

func NewClient (addr string) *Client  {
	return &Client{
		addr: addr,
	}
}

func (client *Client) Get(key string) (string, error) {
   	resp, err := http.Get(client.addr + "?key=" + key)
	if (err != nil) {
		return nil, err
	} else {
		defer resp.Body.Close()
		content, _ := ioutil.ReadAll(resp.Body)
		return content
	}
}

func (client *Client) Set(key string, value string, ttl int) error {
	_, err := http.Post(client.addr + "?key=" + key + "&value=" + value + "ttl=" + ttl, "", nil)
	if (err != nil) {
		return err
	} else {
		return nil
	}

}



