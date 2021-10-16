package accounts

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	tests := []struct {
		name           string
		accountData    *AccountData
		httpUtilsSetup func(*mockHttpUtils)
		wantErr        bool
	}{
		{
			name: "Failed to create an account with a conflict or bad request error",
			httpUtilsSetup: func(client *mockHttpUtils) {
				client.On("Post", mock.Anything, mock.Anything).Return(
					nil,
					errors.New("failed because of an 409, 400"),
				)
			},
			wantErr: true,
		},
		{
			name: "Failed to unsmarshal the response data after creating an account",
			httpUtilsSetup: func(client *mockHttpUtils) {
				client.On("Post", mock.Anything, mock.Anything).Return(
					[]byte("this is an invalid response data"),
					nil,
				)
			},
			wantErr: true,
		},
		{
			name: "Successfully creates an account",
			httpUtilsSetup: func(client *mockHttpUtils) {
				client.On("Post", mock.Anything, mock.Anything).Return(
					loadTestFile("./testdata/fetch_response.json"),
					nil,
				)
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
