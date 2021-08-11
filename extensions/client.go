// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package extensions

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Client struct
type Client struct {
	client *http.Client
}

// NewClient ctor
func NewClient(config *Config) *Client {
	h := &Client{}
	h.client = &http.Client{
		Timeout: 20 * time.Minute,
	}
	return h
}

// Get does a get request
func (c *Client) Get(url, body string) ([]byte, int, error) {
	return c.requestWithBody("GET", url, body)
}

// Put does a put request
func (c *Client) Put(url, body string) ([]byte, int, error) {
	return c.requestWithBody("PUT", url, body)
}

// Post does a post request
func (c *Client) Post(url, body string) ([]byte, int, error) {
	return c.requestWithBody("POST", url, body)
}

// Delete does a put request
func (c *Client) Delete(url string) ([]byte, int, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, 0, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}
	return responseBody, res.StatusCode, nil
}

func (c *Client) requestWithBody(method, url, body string) ([]byte, int, error) {
	ioBody := strings.NewReader(body)

	req, err := http.NewRequest(method, url, ioBody)
	if err != nil {
		return nil, 0, err
	}
	req.Close = true

	res, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}
	return responseBody, res.StatusCode, nil
}
