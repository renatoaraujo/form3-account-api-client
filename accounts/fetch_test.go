package accounts

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFetchAccount(t *testing.T) {
	tests := []struct {
		name           string
		uuid           string
		httpUtilsSetup func(utils *mockHttpUtils)
		wantErr        bool
	}{
		{
			name:    "failed to fetch an account due an invalid uuid",
			uuid:    "invalid-uuid",
			wantErr: true,
		},
		{
			name: "failed to fetch an account due a not found account id",
			uuid: "00000000-0000-0000-0000-000000000000",
			httpUtilsSetup: func(client *mockHttpUtils) {
				client.On("Get", mock.Anything).Return(
					nil,
					errors.New("not found"),
				)
			},
			wantErr: true,
		},
		{
			name: "Failed to fetch data due an data response inconsistency",
			uuid: "00000000-0000-0000-0000-000000000000",
			httpUtilsSetup: func(client *mockHttpUtils) {
				client.On("Get", mock.Anything).Return(
					[]byte("invalid json"),
					errors.New("unable to unmarshal"),
				)
			},
			wantErr: true,
		},
		{
			name: "Successfully fetch an account",
			uuid: "00000000-0000-0000-0000-000000000000",
			httpUtilsSetup: func(client *mockHttpUtils) {
				client.On("Get", mock.Anything).Return(
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
			accountData, err := accountsClient.FetchResource(tt.uuid)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if accountData != nil {
				assert.IsType(t, &AccountData{}, accountData)
			}
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
