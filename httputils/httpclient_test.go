package httputils

import (
	"bytes"
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
		baseURL string
		want    *Client
		wantErr bool
	}{
		{
			name:    "Failed to create client with an invalid base url",
			baseURL: "not-valid-url",
			wantErr: true,
		},
		{
			name:    "Successfully creates new client",
			baseURL: "https://valid-url.com",
			want: &Client{
				httpClient: http.DefaultClient,
				baseURL: url.URL{
					Scheme: "https",
					Host:   "valid-url.com",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(http.DefaultClient, tt.baseURL)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClientPost(t *testing.T) {
	tests := []struct {
		name            string
		httpClientSetup func(*mockHttpClient)
		body            []byte
		want            []byte
		wantErr         bool
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
			name: "Failed perform the post request and receive 409 status code with a valid json data in body",
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
			wantErr: true,
		},
		{
			name: "Failed perform the post request and receive 400 status code with a valid json data in body",
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
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClientMock := &mockHttpClient{}
			if tt.httpClientSetup != nil {
				tt.httpClientSetup(httpClientMock)
			}

			client, err := NewClient(httpClientMock, "http://this-is.fake")
			require.NoError(t, err)

			got, err := client.Post("/a-valid-path", tt.body)
			if tt.wantErr {
				require.Error(t, err)
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
		name            string
		httpClientSetup func(*mockHttpClient)
		want            []byte
		wantErr         bool
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
			wantErr: true,
		},
		{
			name: "Failed perform the get request and receive 400 status code with a valid json data in body",
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
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClientMock := &mockHttpClient{}
			if tt.httpClientSetup != nil {
				tt.httpClientSetup(httpClientMock)
			}

			client, err := NewClient(httpClientMock, "http://this-is.fake")
			require.NoError(t, err)

			got, err := client.Get("/a-valid-path")
			if tt.wantErr {
				require.Error(t, err)
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
		name            string
		httpClientSetup func(*mockHttpClient)
		wantErr         bool
	}{
		{
			name: "Successfully perform the delete request and receive 204 status code with an empty body",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					&http.Response{
						StatusCode: 204,
						Body:       ioutil.NopCloser(bytes.NewBufferString("{}")),
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
						Body:       ioutil.NopCloser(bytes.NewBufferString("{}")),
					},
					nil,
				)
			},
			wantErr: true,
		},
		{
			name: "Failed to perform the delete request and receive 404 status code with an empty body",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					&http.Response{
						StatusCode: 404,
						Body:       ioutil.NopCloser(bytes.NewBufferString("{}")),
					},
					nil,
				)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClientMock := &mockHttpClient{}
			if tt.httpClientSetup != nil {
				tt.httpClientSetup(httpClientMock)
			}

			client, err := NewClient(httpClientMock, "http://this-is.fake")
			require.NoError(t, err)

			query := map[string]string{
				"version": "0",
			}

			err = client.Delete("/a-valid-path", query)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			mock.AssertExpectationsForObjects(t, httpClientMock)
		})
	}
}

func TestHandleResponse(t *testing.T) {
	tests := []struct {
		name       string
		response   *http.Response
		want       []byte
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Successfully return the data with the 2xx status code",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(
					bytes.NewBufferString(
						`{"data": "this is a valid json format"}`,
					),
				),
			},
			want:    []byte(`{"data": "this is a valid json format"}`),
			wantErr: false,
		},
		{
			name: "Failed to return the data with 4xx code",
			response: &http.Response{
				StatusCode: 400,
				Body: ioutil.NopCloser(
					bytes.NewBufferString(
						`{"error_message": "it failed"}`,
					),
				),
			},
			wantErr:    true,
			wantErrMsg: "api failure with status code 400 and message: it failed",
		},
		{
			name: "Failed to return the data with 5xx code and empty body",
			response: &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			},
			wantErr:    true,
			wantErrMsg: "api failure with status code 500 and no message received",
		},
		{
			name: "Failed to return the data with 1xx code and empty body",
			response: &http.Response{
				StatusCode: 100,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			},
			wantErr:    true,
			wantErrMsg: "unexpected status code 100",
		},
		{
			name: "Failed to return the data with 3xx code and empty body",
			response: &http.Response{
				StatusCode: 301,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			},
			wantErr:    true,
			wantErrMsg: "unexpected status code 301",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handleResponse(tt.response)
			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.wantErrMsg, err.Error())
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
