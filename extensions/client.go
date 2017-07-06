// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package extensions

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Client struct
type Client struct {
	config *Config
	client *http.Client
}

// NewClient ctor
func NewClient(config *Config) *Client {
	h := &Client{
		config: config,
	}
	h.client = &http.Client{
		Timeout: 20 * time.Minute,
	}
	return h
}

// Get does a get request
func (c *Client) Get(url string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}
	c.addAuthHeader(req)
	res, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}
	return body, res.StatusCode, nil
}

// Put does a put request
func (c *Client) Put(url string, body string) ([]byte, int, error) {
	return c.putOrPost("PUT", url, body)
}

// Post does a post request
func (c *Client) Post(url string, body string) ([]byte, int, error) {
	return c.putOrPost("POST", url, body)
}

// Delete does a put request
func (c *Client) Delete(url string) ([]byte, int, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, 0, err
	}

	c.addAuthHeader(req)
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

func (c *Client) putOrPost(method, url string, body string) ([]byte, int, error) {
	ioBody := strings.NewReader(body)

	req, err := http.NewRequest(method, url, ioBody)
	if err != nil {
		return nil, 0, err
	}
	req.Close = true

	c.addAuthHeader(req)
	res, err := c.client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("Error creating cluster")
	}
	defer res.Body.Close()
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}
	return responseBody, res.StatusCode, nil
}

func (c *Client) addAuthHeader(req *http.Request) {
	if c.config != nil {
		auth := fmt.Sprintf("Bearer %s", c.config.Token)
		req.Header.Add("Authorization", auth)
	}

	req.Header.Add("Content-Type", "application/json")
}
