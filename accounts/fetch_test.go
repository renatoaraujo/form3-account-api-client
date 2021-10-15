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
		AccountID      string
		httpUtilsSetup func(utils *mockHttpUtils)
		wantErr        bool
	}{
		{
			name:      "Failed to fetch an account due a not found account id",
			AccountID: "00000000-0000-0000-0000-000000000000",
			httpUtilsSetup: func(client *mockHttpUtils) {
				client.On("Get", mock.Anything).Return(
					nil,
					errors.New("not found"),
				)
			},
			wantErr: true,
		},
		{
			name:      "Failed to fetch data due a missing data in response format",
			AccountID: "00000000-0000-0000-0000-000000000000",
			httpUtilsSetup: func(client *mockHttpUtils) {
				client.On("Get", mock.Anything).Return(
					[]byte("invalid json"),
					errors.New("unable to unmarshal invalid json"),
				)
			},
			wantErr: true,
		},
		{
			name:      "Successfully fetch an account",
			AccountID: "00000000-0000-0000-0000-000000000000",
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
			accountData, err := accountsClient.FetchResource(tt.AccountID)
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
