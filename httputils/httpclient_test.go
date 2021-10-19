package httputils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	tests := []struct {
		name    string
		baseURI string
		timeout int
		wantErr bool
	}{
		{
			name:    "Failed to create client with an invalid base url",
			baseURI: "not-valid-url",
			wantErr: true,
		},
		{
			name:    "Successfully creates new client",
			baseURI: "https://valid-url.com",
			timeout: 15,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.baseURI, tt.timeout)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.IsType(t, &Client{}, got)
		})
	}
}

func TestClientPost(t *testing.T) {
	tests := []struct {
		name             string
		httpClientSetup  func(*mockHttpClient)
		bodyReader       func(io.Reader) ([]byte, error)
		respUnmarshaller func([]byte, interface{}) error
		reqCreator       func(method, url string, body io.Reader) (*http.Request, error)
		want             []byte
		wantErr          bool
		wantErrMsg       string
	}{
		{
			name: "Successfully perform the post request and receive 201 status code with a valid json data in body",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					&http.Response{
						StatusCode: 201,
						Body: ioutil.NopCloser(
							bytes.NewBufferString(
								`{"data":"some valid json data"}`,
							),
						),
					},
					nil,
				)
			},
			want:    []byte(`{"data":"some valid json data"}`),
			wantErr: false,
		},
		{
			name: "Failed to perform the post request and receive 409 status code with a valid json data in body",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					&http.Response{
						StatusCode: 409,
						Body: ioutil.NopCloser(
							bytes.NewBufferString(
								`{"error_message":"it violates a duplicate constraint"}`,
							),
						),
					},
					nil,
				)
			},
			wantErr:    true,
			wantErrMsg: "api failure with status code 409 and message: it violates a duplicate constraint",
		},
		{
			name: "Failed to perform the post request and receive 400 status code with a valid json data in body",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					&http.Response{
						StatusCode: 400,
						Body: ioutil.NopCloser(
							bytes.NewBufferString(
								`{"error_message":"validation failure"}`,
							),
						),
					},
					nil,
				)
			},
			wantErr:    true,
			wantErrMsg: "api failure with status code 400 and message: validation failure",
		},
		{
			name: "Failed to perform the post request and receive 500 status code with an empty",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					&http.Response{
						StatusCode: 500,
						Body:       ioutil.NopCloser(bytes.NewBufferString("")),
					},
					nil,
				)
			},
			wantErr:    true,
			wantErrMsg: "unexpected status code 500",
		},
		{
			name: "Failed to perform the request failing the http client",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					nil,
					errors.New("failed to perform request"),
				)
			},
			wantErr:    true,
			wantErrMsg: "failed to perform request; failed to post data",
		},
		{
			name: "Failed to create the request",
			reqCreator: func(string, string, io.Reader) (*http.Request, error) {
				return nil, errors.New("failed to create the request")
			},
			wantErr:    true,
			wantErrMsg: "failed to create the request",
		},
		{
			name: "Failed to read the response body",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					&http.Response{
						StatusCode: 201,
						Body: ioutil.NopCloser(
							bytes.NewBufferString(
								`{"data":"some valid data"}`,
							),
						),
					},
					nil,
				)
			},
			bodyReader: func(io.Reader) ([]byte, error) {
				return nil, errors.New("failed to read body for some reason")
			},
			wantErr:    true,
			wantErrMsg: "failed to read body for some reason; failed to read response body",
		},
		{
			name: "Failed to convert the error response body",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					&http.Response{
						StatusCode: 400,
						Body: ioutil.NopCloser(
							bytes.NewBufferString(
								`{"error":"this is not the structure expected"}`,
							),
						),
					},
					nil,
				)
			},
			respUnmarshaller: func([]byte, interface{}) error {
				return errors.New("failed to unmarshal")
			},
			wantErr:    true,
			wantErrMsg: "failed to unmarshal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClientMock := &mockHttpClient{}
			if tt.httpClientSetup != nil {
				tt.httpClientSetup(httpClientMock)
			}
			client := createFakeHttpClient(httpClientMock, tt.bodyReader, tt.respUnmarshaller, tt.reqCreator)

			got, err := client.Post("/a-valid-path", []byte("something"))
			if tt.wantErr {
				require.Error(t, err)
				assert.EqualError(t, err, tt.wantErrMsg)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
			mock.AssertExpectationsForObjects(t, httpClientMock)
		})
	}
}

func TestClientGet(t *testing.T) {
	tests := []struct {
		name             string
		httpClientSetup  func(*mockHttpClient)
		bodyReader       func(io.Reader) ([]byte, error)
		respUnmarshaller func([]byte, interface{}) error
		reqCreator       func(method, url string, body io.Reader) (*http.Request, error)
		want             []byte
		wantErr          bool
		wantErrMsg       string
	}{
		{
			name: "Successfully perform the get request and receive 200 status code with a valid json data in body",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					&http.Response{
						StatusCode: 200,
						Body: ioutil.NopCloser(
							bytes.NewBufferString(
								`{"data":"this is a valid json data"}`,
							),
						),
					},
					nil,
				)
			},
			want:    []byte(`{"data":"this is a valid json data"}`),
			wantErr: false,
		},
		{
			name: "Failed to perform the get request and receive 404 status code with a valid json data in body",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					&http.Response{
						StatusCode: 404,
						Body: ioutil.NopCloser(
							bytes.NewBufferString(
								`{"error_message":"record xxx-xxx does not exist"}`,
							),
						),
					},
					nil,
				)
			},
			wantErr:    true,
			wantErrMsg: "api failure with status code 404 and message: record xxx-xxx does not exist",
		},
		{
			name: "Failed to perform the get request and receive 400 status code with a valid json data in body",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					&http.Response{
						StatusCode: 400,
						Body: ioutil.NopCloser(
							bytes.NewBufferString(
								`{"error_message":"id is not a valid uuid"}`,
							),
						),
					},
					nil,
				)
			},
			wantErr:    true,
			wantErrMsg: "api failure with status code 400 and message: id is not a valid uuid",
		},
		{
			name: "Failed to perform the get request and receive 500 status code with an empty body",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					&http.Response{
						StatusCode: 500,
						Body:       ioutil.NopCloser(bytes.NewBufferString("")),
					},
					nil,
				)
			},
			wantErr:    true,
			wantErrMsg: "unexpected status code 500",
		},
		{
			name: "Failed to perform the request failing the http client",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					nil,
					errors.New("failed to perform request"),
				)
			},
			wantErr:    true,
			wantErrMsg: "failed to perform request",
		},
		{
			name: "Failed to create the request",
			reqCreator: func(method, url string, body io.Reader) (*http.Request, error) {
				return nil, errors.New("failed to create the request")
			},
			wantErr:    true,
			wantErrMsg: "failed to create the request",
		},
		{
			name: "Failed to convert the error response body",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					&http.Response{
						StatusCode: 400,
						Body: ioutil.NopCloser(
							bytes.NewBufferString(
								`{"error":"this is not the structure expected"}`,
							),
						),
					},
					nil,
				)
			},
			respUnmarshaller: func([]byte, interface{}) error {
				return errors.New("failed to unmarshal")
			},
			wantErr:    true,
			wantErrMsg: "failed to unmarshal",
		},
		{
			name: "Failed to read the response body",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					&http.Response{
						StatusCode: 200,
						Body: ioutil.NopCloser(
							bytes.NewBufferString(
								`{"data":"some valid data"}`,
							),
						),
					},
					nil,
				)
			},
			bodyReader: func(io.Reader) ([]byte, error) {
				return nil, errors.New("failed to read body for some reason")
			},
			wantErr:    true,
			wantErrMsg: "failed to read body for some reason; failed to read response body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClientMock := &mockHttpClient{}
			if tt.httpClientSetup != nil {
				tt.httpClientSetup(httpClientMock)
			}

			client := createFakeHttpClient(httpClientMock, tt.bodyReader, tt.respUnmarshaller, tt.reqCreator)

			got, err := client.Get("/a-valid-path")
			if tt.wantErr {
				require.Error(t, err)
				assert.EqualError(t, err, tt.wantErrMsg)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
			mock.AssertExpectationsForObjects(t, httpClientMock)
		})
	}
}

func TestClientDelete(t *testing.T) {
	tests := []struct {
		name             string
		httpClientSetup  func(*mockHttpClient)
		bodyReader       func(io.Reader) ([]byte, error)
		respUnmarshaller func([]byte, interface{}) error
		reqCreator       func(method, url string, body io.Reader) (*http.Request, error)
		wantErr          bool
		wantErrMsg       string
	}{
		{
			name: "Successfully perform the delete request and receive 204 status code with an empty body",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					&http.Response{
						StatusCode: 204,
						Body:       ioutil.NopCloser(bytes.NewBufferString("")),
					},
					nil,
				)
			},
			wantErr: false,
		},
		{
			name: "Failed to perform the delete request and receive 400 status code with an empty body",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					&http.Response{
						StatusCode: 400,
						Body: ioutil.NopCloser(
							bytes.NewBufferString(`{"error_message":"invalid version number"}`),
						),
					},
					nil,
				)
			},
			wantErr:    true,
			wantErrMsg: "api failure with status code 400 and message: invalid version number",
		},
		{
			name: "Failed to perform the delete request and receive 404 status code with an empty body",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					&http.Response{
						StatusCode: 404,
						Body:       ioutil.NopCloser(bytes.NewBufferString("")),
					},
					nil,
				)
			},
			wantErr:    true,
			wantErrMsg: "api failure with status code 404 and message: not found",
		},
		{
			name: "Failed to perform the delete request and receive 500 status code with an empty body",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					&http.Response{
						StatusCode: 500,
						Body: ioutil.NopCloser(
							bytes.NewBufferString(""),
						),
					},
					nil,
				)
			},
			wantErr:    true,
			wantErrMsg: "unexpected status code 500",
		},
		{
			name: "Failed to perform the request failing the http client",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					nil,
					errors.New("failed to perform request"),
				)
			},
			wantErr:    true,
			wantErrMsg: "failed to perform request",
		},
		{
			name: "Failed to create the request",
			reqCreator: func(string, string, io.Reader) (*http.Request, error) {
				return nil, errors.New("failed to create the request")
			},
			wantErr:    true,
			wantErrMsg: "failed to create the request",
		},
		{
			name: "Failed to read the error response body",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					&http.Response{
						StatusCode: 400,
						Body: ioutil.NopCloser(
							bytes.NewBufferString(
								`{"error":"this is not the expected format"}`,
							),
						),
					},
					nil,
				)
			},
			bodyReader: func(io.Reader) ([]byte, error) {
				return nil, errors.New("failed to read body for some reason")
			},
			wantErr:    true,
			wantErrMsg: "failed to read body for some reason; failed to read response body",
		},
		{
			name: "Failed to convert the error response body",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					&http.Response{
						StatusCode: 400,
						Body: ioutil.NopCloser(
							bytes.NewBufferString(
								`{"error":"this is not the structure expected"}`,
							),
						),
					},
					nil,
				)
			},
			respUnmarshaller: func([]byte, interface{}) error {
				return errors.New("failed to unmarshal")
			},
			wantErr:    true,
			wantErrMsg: "failed to unmarshal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClientMock := &mockHttpClient{}
			if tt.httpClientSetup != nil {
				tt.httpClientSetup(httpClientMock)
			}
			client := createFakeHttpClient(httpClientMock, tt.bodyReader, tt.respUnmarshaller, tt.reqCreator)

			query := map[string]string{
				"version": "0",
			}

			err := client.Delete("/a-valid-path", query)
			if tt.wantErr {
				require.Error(t, err)
				assert.EqualError(t, err, tt.wantErrMsg)
			} else {
				require.NoError(t, err)
			}

			mock.AssertExpectationsForObjects(t, httpClientMock)
		})
	}
}

func createFakeHttpClient(
	mock *mockHttpClient,
	bodyReader func(io.Reader) ([]byte, error),
	respUnmarshaller func([]byte, interface{}) error,
	reqCreator func(method, url string, body io.Reader) (*http.Request, error),
) Client {
	if bodyReader == nil {
		bodyReader = ioutil.ReadAll
	}

	if respUnmarshaller == nil {
		respUnmarshaller = json.Unmarshal
	}

	if reqCreator == nil {
		reqCreator = http.NewRequest
	}

	return Client{
		httpClient: mock,
		baseURI: url.URL{
			Scheme: "https",
			Host:   "api.form3.tech",
		},
		bodyReader:       bodyReader,
		respUnmarshaller: respUnmarshaller,
		reqCreator:       reqCreator,
	}
}
