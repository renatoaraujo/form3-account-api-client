package accounts

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeleteAccount(t *testing.T) {
	tests := []struct {
		name            string
		uuid            string
		httpClientSetup func(*mockHttpClient)
		wantErr         bool
	}{
		{
			name:    "failed to delete an account due an invalid uuid",
			uuid:    "invalid-uuid",
			wantErr: true,
		},
		{
			name: "failed to delete an account due a response with an error content from the api",
			uuid: "00000000-0000-0000-0000-000000000000",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Delete", mock.Anything).Return(errors.New("failed because of an 404 or 409"))
			},
			wantErr: true,
		},
		{
			name: "it will successfully delete an account",
			uuid: "00000000-0000-0000-0000-000000000000",
			httpClientSetup: func(client *mockHttpClient) {
				client.On("Delete", mock.Anything).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			httpClientMock := &mockHttpClient{}
			if tt.httpClientSetup != nil {
				tt.httpClientSetup(httpClientMock)
			}

			accountsClient := NewClient(httpClientMock)
			err := accountsClient.DeleteResource(tt.uuid, 123)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
