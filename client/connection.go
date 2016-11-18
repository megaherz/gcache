package client

import (
	"hash/crc32"
	"io"
	"net/http"
)

type Connection struct{
	addr string
	psw string
}

type Connections []Connection

func (c Connections) getNode(key string) Connection {
	hash := crc32.ChecksumIEEE([]byte(key))
	n := hash % uint32(len(c))
	return c[n]
}

func (conn *Connection) doRequest(method, urlStr string, body io.Reader) (*http.Response, error) {

	req, err := http.NewRequest(method, conn.addr + urlStr, body)
	if err != nil {
		return nil, err
	}

	// Set the authorization header if password is set
	if conn.psw != "" {
		req.Header.Set(headerAuthorization, conn.psw)
	}

	return http.DefaultClient.Do(req)
}

type httpResponse struct {
	url      string
	response *http.Response
	err      error
}

func (conns Connections) doParallelGetRequest(query string) []*httpResponse {
	ch := make(chan *httpResponse)
	responses := []*httpResponse{}

	for _, conn := range conns {
		go func(url string) {
			resp, err := conn.doRequest(http.MethodGet, url, nil)
			ch <- &httpResponse{url, resp, err}
			if err != nil && resp != nil && resp.StatusCode == http.StatusOK {
				resp.Body.Close()
			}
		}(query)
	}

	for {
		select {
		case r := <-ch:
			responses = append(responses, r)
			if len(responses) == len(conns) {
				return responses
			}
		}
	}
	return responses
}

