package httputils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	httpClient httpClient
	baseURL    url.URL
}

func NewClient(client httpClient, baseURI string) (*Client, error) {
	parsedBaseURI, err := url.ParseRequestURI(baseURI)
	if err != nil {
		return nil, fmt.Errorf("%w; invalid base uri", err)
	}

	return &Client{
		httpClient: client,
		baseURL: url.URL{
			Scheme: parsedBaseURI.Scheme,
			Host:   parsedBaseURI.Host,
		},
	}, nil
}

func (c Client) Post(resourcePath string, body []byte) ([]byte, error) {
	requestURL := c.baseURL.ResolveReference(&url.URL{Path: resourcePath})
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

func (c Client) Get(resourcePath string) ([]byte, error) {
	requestURL := c.baseURL.ResolveReference(&url.URL{Path: resourcePath})
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

func (c Client) Delete(resourcePath string, query map[string]string) error {
	rawQuery := url.Values{}
	for key, value := range query {
		rawQuery.Add(key, value)
	}
	requestURL := c.baseURL.ResolveReference(&url.URL{Path: resourcePath, RawQuery: rawQuery.Encode()})
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
