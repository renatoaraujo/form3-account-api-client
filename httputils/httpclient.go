package httputils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	httpClient       httpClient
	baseURI          url.URL
	bodyReader       bodyReader
	respUnmarshaller respUnmarshaller
	reqCreator       reqCreator
}

type bodyReader func(io.Reader) ([]byte, error)
type respUnmarshaller func([]byte, interface{}) error
type reqCreator func(method, url string, body io.Reader) (*http.Request, error)

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
		baseURI: url.URL{
			Scheme: parsedBaseURI.Scheme,
			Host:   parsedBaseURI.Host,
		},
		bodyReader:       ioutil.ReadAll,
		respUnmarshaller: json.Unmarshal,
		reqCreator:       http.NewRequest,
	}, nil
}

// Post data to an API endpoint with given path and body content
func (c Client) Post(resourcePath string, body []byte) ([]byte, error) {
	requestURL := c.baseURI.ResolveReference(&url.URL{Path: resourcePath})
	request, err := c.reqCreator(http.MethodPost, requestURL.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%w; failed to post data", err)
	}
	defer response.Body.Close()

	respBody, err := c.bodyReader(response.Body)
	if err != nil {
		return nil, fmt.Errorf("%w; failed to read response body", err)
	}

	switch response.StatusCode {
	case http.StatusCreated:
		return respBody, nil
	case http.StatusConflict, http.StatusBadRequest:
		var errRes ResponseError
		if err := c.respUnmarshaller(respBody, &errRes); err != nil {
			return nil, err
		}

		errRes.StatusCode = response.StatusCode
		return nil, &errRes
	default:
		return nil, fmt.Errorf("unexpected status code %d", response.StatusCode)
	}
}

// Get data from an API endpoint with given path
func (c Client) Get(resourcePath string) ([]byte, error) {
	requestURL := c.baseURI.ResolveReference(&url.URL{Path: resourcePath})
	request, err := c.reqCreator(http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	respBody, err := c.bodyReader(response.Body)
	if err != nil {
		return nil, fmt.Errorf("%w; failed to read response body", err)
	}

	switch response.StatusCode {
	case http.StatusOK:
		return respBody, nil
	case http.StatusNotFound, http.StatusBadRequest:
		var errRes ResponseError
		if err := c.respUnmarshaller(respBody, &errRes); err != nil {
			return nil, err
		}

		errRes.StatusCode = response.StatusCode
		return nil, &errRes
	default:
		return nil, fmt.Errorf("unexpected status code %d", response.StatusCode)
	}
}

// Delete data from an API endpoint with given path and query string
func (c Client) Delete(resourcePath string, query map[string]string) error {
	rawQuery := url.Values{}
	for key, value := range query {
		rawQuery.Add(key, value)
	}
	requestURL := c.baseURI.ResolveReference(&url.URL{Path: resourcePath, RawQuery: rawQuery.Encode()})
	request, err := c.reqCreator(http.MethodDelete, requestURL.String(), nil)
	if err != nil {
		return err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		respBody, err := c.bodyReader(response.Body)
		if err != nil {
			return fmt.Errorf("%w; failed to read response body", err)
		}
		var errRes ResponseError
		if err := c.respUnmarshaller(respBody, &errRes); err != nil {
			return err
		}

		errRes.StatusCode = response.StatusCode
		return &errRes
	case http.StatusNotFound:
		return &ResponseError{
			ErrorMessage: "not found",
			StatusCode:   404,
		}
	default:
		return fmt.Errorf("unexpected status code %d", response.StatusCode)
	}
}
