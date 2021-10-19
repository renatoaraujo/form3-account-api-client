package accounts

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateResource(t *testing.T) {
	tests := []struct {
		name              string
		accountData       *AccountData
		httpUtilsSetup    func(*mockHttpUtils)
		respUnmarshaller  func([]byte, interface{}) error
		payloadMarshaller func(v interface{}) ([]byte, error)
		wantErr           bool
	}{
		{
			name: "Failed to create an account because of an API error",
			httpUtilsSetup: func(client *mockHttpUtils) {
				client.On("Post", mock.Anything, mock.Anything).Return(
					nil,
					errors.New("the api failed the request"),
				)
			},
			wantErr: true,
		},
		{
			name: "Failed to convert the response data after creating an account successfully",
			httpUtilsSetup: func(client *mockHttpUtils) {
				client.On("Post", mock.Anything, mock.Anything).Return(
					[]byte("the api did not failed but this is a wrong response data format"),
					nil,
				)
			},
			wantErr: true,
		},
		{
			name: "Successfully creates an account",
			httpUtilsSetup: func(client *mockHttpUtils) {
				client.On("Post", mock.Anything, mock.Anything).Return(
					loadTestFile("./testdata/api_response.json"),
					nil,
				)
			},
			wantErr: false,
		},
		{
			name: "Failed to marshal the payload",
			payloadMarshaller: func(interface{}) ([]byte, error) {
				return nil, errors.New("failed to marshal")
			},
			wantErr: true,
		},
		{
			name: "Failed to unmarshal the successful response",
			httpUtilsSetup: func(client *mockHttpUtils) {
				client.On("Post", mock.Anything, mock.Anything).Return(
					loadTestFile("./testdata/api_response.json"),
					nil,
				)
			},
			respUnmarshaller: func([]byte, interface{}) error {
				return errors.New("failed to unmarshal")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			httpUtilsMock := &mockHttpUtils{}
			if tt.httpUtilsSetup != nil {
				tt.httpUtilsSetup(httpUtilsMock)
			}

			if tt.respUnmarshaller == nil {
				tt.respUnmarshaller = json.Unmarshal
			}

			if tt.payloadMarshaller == nil {
				tt.payloadMarshaller = json.Marshal
			}

			accountsClient := Client{
				http:              httpUtilsMock,
				respUnmarshaller:  tt.respUnmarshaller,
				payloadMarshaller: tt.payloadMarshaller,
			}
			accountData, err := accountsClient.CreateResource(tt.accountData)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if accountData != nil {
				assert.IsType(t, &AccountData{}, accountData)
			}

			mock.AssertExpectationsForObjects(t, httpUtilsMock)
		})
	}
}

func TestFetchResource(t *testing.T) {
	tests := []struct {
		name             string
		httpUtilsSetup   func(utils *mockHttpUtils)
		respUnmarshaller func([]byte, interface{}) error
		wantErr          bool
	}{
		{
			name: "Failed to fetch account data because of account id was not found",
			httpUtilsSetup: func(client *mockHttpUtils) {
				client.On("Get", mock.Anything).Return(
					nil,
					errors.New("not found"),
				)
			},
			wantErr: true,
		},
		{
			name: "Failed to fetch because of an invalid format from the api response",
			httpUtilsSetup: func(client *mockHttpUtils) {
				client.On("Get", mock.Anything).Return(
					[]byte("invalid json"),
					errors.New("unable to unmarshal invalid json"),
				)
			},
			wantErr: true,
		},
		{
			name: "Successfully fetches an account",
			httpUtilsSetup: func(client *mockHttpUtils) {
				client.On("Get", mock.Anything).Return(
					loadTestFile("./testdata/api_response.json"),
					nil,
				)
			},
			wantErr: false,
		},
		{
			name: "Failed to unmarshal the successful response",
			httpUtilsSetup: func(client *mockHttpUtils) {
				client.On("Get", mock.Anything).Return(
					loadTestFile("./testdata/api_response.json"),
					nil,
				)
			},
			respUnmarshaller: func([]byte, interface{}) error {
				return errors.New("failed to unmarshal")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			httpUtilsMock := &mockHttpUtils{}
			if tt.httpUtilsSetup != nil {
				tt.httpUtilsSetup(httpUtilsMock)
			}

			if tt.respUnmarshaller == nil {
				tt.respUnmarshaller = json.Unmarshal
			}

			accountsClient := Client{
				http:              httpUtilsMock,
				respUnmarshaller:  tt.respUnmarshaller,
				payloadMarshaller: json.Marshal,
			}

			accountID, err := uuid.NewUUID()
			require.NoError(t, err)

			accountData, err := accountsClient.FetchResource(accountID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if accountData != nil {
				assert.IsType(t, &AccountData{}, accountData)
			}

			mock.AssertExpectationsForObjects(t, httpUtilsMock)
		})
	}
}

func TestDeleteResource(t *testing.T) {
	tests := []struct {
		name           string
		httpUtilsSetup func(*mockHttpUtils)
		wantErr        bool
	}{
		{
			name: "Failed to delete an account with an error response from the api",
			httpUtilsSetup: func(client *mockHttpUtils) {
				client.On("Delete", mock.Anything, mock.Anything).Return(
					errors.New("failed because of a failure in the api"),
				)
			},
			wantErr: true,
		},
		{
			name: "Successfully deletes an account",
			httpUtilsSetup: func(client *mockHttpUtils) {
				client.On("Delete", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			httpUtilsMock := &mockHttpUtils{}
			if tt.httpUtilsSetup != nil {
				tt.httpUtilsSetup(httpUtilsMock)
			}

			accountsClient := NewClient(httpUtilsMock)

			accountID, err := uuid.NewUUID()
			require.NoError(t, err)

			err = accountsClient.DeleteResource(accountID, 123)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			mock.AssertExpectationsForObjects(t, httpUtilsMock)
		})
	}
}

func loadTestFile(file string) []byte {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		panic("failed to load the test data file")
	}

	return raw
}
