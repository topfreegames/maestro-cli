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
func (c *Client) Get(url, body string, headers ...map[string]string) ([]byte, int, error) {
	var requestHeaders map[string]string
	if len(headers) > 0 {
		requestHeaders = headers[0]
	}
	return c.requestWithBody("GET", url, body, requestHeaders)

	// req, err := http.NewRequest("GET", url, nil)
	// if err != nil {
	// 	return nil, 0, err
	// }
	// var requestHeaders map[string]string
	// if len(headers) > 0 {
	// 	requestHeaders = headers[0]
	// }
	// c.addHeaders(req, requestHeaders)
	// res, err := c.client.Do(req)
	// if err != nil {
	// 	return nil, 0, err
	// }
	// defer res.Body.Close()
	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	return nil, 0, err
	// }
	// return body, res.StatusCode, nil
}

// Put does a put request
func (c *Client) Put(url, body string, headers ...map[string]string) ([]byte, int, error) {
	var requestHeaders map[string]string
	if len(headers) > 0 {
		requestHeaders = headers[0]
	}
	return c.requestWithBody("PUT", url, body, requestHeaders)
}

// Post does a post request
func (c *Client) Post(url, body string, headers ...map[string]string) ([]byte, int, error) {
	var requestHeaders map[string]string
	if len(headers) > 0 {
		requestHeaders = headers[0]
	}
	return c.requestWithBody("POST", url, body, requestHeaders)
}

// Delete does a put request
func (c *Client) Delete(url string, headers ...map[string]string) ([]byte, int, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, 0, err
	}

	var requestHeaders map[string]string
	if len(headers) > 0 {
		requestHeaders = headers[0]
	}
	c.addHeaders(req, requestHeaders)
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

func (c *Client) requestWithBody(method, url, body string, headers map[string]string) ([]byte, int, error) {
	ioBody := strings.NewReader(body)

	req, err := http.NewRequest(method, url, ioBody)
	if err != nil {
		return nil, 0, err
	}
	req.Close = true

	c.addHeaders(req, headers)
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

func (c *Client) addHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	c.addAuthHeader(req)
}

func (c *Client) addAuthHeader(req *http.Request) {
	if c.config != nil {
		auth := fmt.Sprintf("Bearer %s", c.config.Token)
		req.Header.Add("Authorization", auth)
	}

	req.Header.Add("Content-Type", "application/json")
}
