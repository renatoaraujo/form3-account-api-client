package httputils

import (
	"bytes"
	"encoding/json"
	"errors"
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

func (c Client) Post(resourcePath string, payload []byte) ([]byte, error) {
	requestURL := c.baseURL.ResolveReference(&url.URL{Path: resourcePath})
	request, err := http.NewRequest("POST", requestURL.String(), bytes.NewBuffer(payload))
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

func handleResponse(response *http.Response) ([]byte, error) {
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("%w; failed to read response body", err)
	}

	if !json.Valid(respBody) {
		return nil, errors.New("invalid json format in the response body")
	}

	switch response.StatusCode {
	case 200, 201:
		return respBody, nil
	case 400, 404, 409:
		respError := &responseError{}
		err = json.Unmarshal(respBody, respError)
		if err != nil || respError.ErrorMessage == "" {
			return nil, fmt.Errorf("%w; failed to unmarshal response data", err)
		}

		return nil, errors.New(respError.ErrorMessage)
	default:
		return nil, fmt.Errorf("unexpected status code %d", response.StatusCode)
	}
}
