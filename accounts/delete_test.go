package accounts

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeleteAccount(t *testing.T) {
	tests := []struct {
		name           string
		accountID      string
		httpUtilsSetup func(*mockHttpUtils)
		wantErr        bool
	}{
		{
			name:      "failed to delete an account due a response with an error content from the api",
			accountID: "00000000-0000-0000-0000-000000000000",
			httpUtilsSetup: func(client *mockHttpUtils) {
				client.On("Delete", mock.Anything).Return(
					errors.New("failed because of an 404 or 409"),
				)
			},
			wantErr: true,
		},
		{
			name:      "it will successfully delete an account",
			accountID: "00000000-0000-0000-0000-000000000000",
			httpUtilsSetup: func(client *mockHttpUtils) {
				client.On("Delete", mock.Anything).Return(nil)
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
			err := accountsClient.DeleteResource(tt.accountID, 123)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
