package httputils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		wantErr bool
	}{
		{
			name:    "Failed to create client with an invalid base url",
			baseURL: "not-valid-url",
			wantErr: true,
		},
		{
			name:    "Successfully creates new client",
			baseURL: "http://valid-url.com",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(&http.Client{
				Timeout: time.Duration(30) * time.Second,
			}, tt.baseURL)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
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
			name: "Failed when trying to post to api with conflict response (status code 409)",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Do", mock.Anything).Return(
					&http.Response{
						StatusCode: 409,
						Body: ioutil.NopCloser(
							bytes.NewBufferString(
								`{"error_message":"Account cannot be created as it violates a duplicate constraint"}`,
							),
						),
					},
					nil,
				)
			},
			wantErr: true,
		},
		{
			name: "Failed when trying to post to api with bad request response (status code 400)",
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
		{
			name: "Success when posting to api with created response (status code 201)",
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

func TestHandleResponse(t *testing.T) {
	tests := []struct {
		name     string
		response *http.Response
		want     []byte
		wantErr  bool
	}{
		{
			name: "Failed due invalid json format in the response body",
			response: &http.Response{
				StatusCode: 201,
				Body: ioutil.NopCloser(
					bytes.NewBufferString(
						`this is an invalid json`,
					),
				),
			},
			wantErr: true,
		},
		{
			name: "Successfully return the data with the status code create (status code 201)",
			response: &http.Response{
				StatusCode: 201,
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
			name: "Failed to return data with status code conflict (status code 409)",
			response: &http.Response{
				StatusCode: 409,
				Body: ioutil.NopCloser(
					bytes.NewBufferString(
						`{"error_message": "it failed because of conflict"}`,
					),
				),
			},
			wantErr: true,
		},
		{
			name: "Failed to return data with status code bad request (status code 400)",
			response: &http.Response{
				StatusCode: 400,
				Body: ioutil.NopCloser(
					bytes.NewBufferString(
						`{"error_message": "it failed because of a bad request (maybe validation)"}`,
					),
				),
			},
			wantErr: true,
		},
		{
			name: "Failed to return data with status code bad request with a empty body (status code 400)",
			response: &http.Response{
				StatusCode: 400,
				Body: ioutil.NopCloser(
					bytes.NewBufferString(
						`{"data": "this should not be data"}`,
					),
				),
			},
			wantErr: true,
		},
		{
			name: "Failed to return data with status i am a tea pot (status code 418)",
			response: &http.Response{
				StatusCode: 418,
				Body: ioutil.NopCloser(
					bytes.NewBufferString(
						`{"data": "I am a tea pot"}`,
					),
				),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handleResponse(tt.response)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
