package httputils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is the representation of the client to perform some http operations
type Client struct {
	httpClient httpClient
	baseURI    *url.URL
}

// NewClient creates a new http client with the base URI and the timeout for the requests made by this client
func NewClient(baseURI string, timeout int) (*Client, error) {
	parsedBaseURI, err := url.ParseRequestURI(baseURI)
	if err != nil {
		return nil, fmt.Errorf("%w; invalid base uri", err)
	}

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	return &Client{
		httpClient: client,
		baseURI:    parsedBaseURI,
	}, nil
}

// Post data to an API endpoint with given path and body content
func (c Client) Post(resourcePath string, body []byte) ([]byte, error) {
	requestURL := c.baseURI.ResolveReference(&url.URL{Path: resourcePath})
	request, err := http.NewRequest("POST", requestURL.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	return handleResponse(response)
}

// Get data from an API endpoint with given path
func (c Client) Get(resourcePath string) ([]byte, error) {
	requestURL := c.baseURI.ResolveReference(&url.URL{Path: resourcePath})
	request, err := http.NewRequest("GET", requestURL.String(), nil)
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	return handleResponse(response)
}

// Delete data from an API endpoint with given path and query string
func (c Client) Delete(resourcePath string, query map[string]string) error {
	rawQuery := url.Values{}
	for key, value := range query {
		rawQuery.Add(key, value)
	}
	requestURL := c.baseURI.ResolveReference(&url.URL{Path: resourcePath, RawQuery: rawQuery.Encode()})
	request, err := http.NewRequest("DELETE", requestURL.String(), nil)
	if err != nil {
		return err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	_, err = handleResponse(response)
	if err != nil {
		return err
	}

	return nil
}

func handleResponse(response *http.Response) ([]byte, error) {
	if response.StatusCode >= 200 && response.StatusCode < 299 {
		respBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("%w; failed to read response body", err)
		}

		return respBody, nil
	}

	if response.StatusCode >= 400 {
		var errRes ResponseError
		_ = json.NewDecoder(response.Body).Decode(&errRes)

		if errRes.StatusCode == 0 {
			errRes.StatusCode = response.StatusCode
		}
		return nil, &errRes
	}

	return nil, fmt.Errorf("unexpected status code %d", response.StatusCode)
}
