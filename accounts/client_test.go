package accounts

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestValidateAccountIDFormat(t *testing.T) {
	tests := []struct {
		name      string
		accountID string
		wantErr   bool
	}{
		{
			name:      "Failed due an invalid account id format",
			accountID: "this is an invalid uuid format",
			wantErr:   true,
		},
		{
			name:      "Successfully validate the account id format",
			accountID: "00000000-0000-0000-0000-000000000000",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAccountIDFormat(tt.accountID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestExtractAccountDataFromResponse(t *testing.T) {
	tests := []struct {
		name        string
		response    []byte
		accountData AccountData
		wantErr     bool
	}{
		{
			name:     "Failed to unsmarshal an invalid response json data format",
			response: []byte("this is an invalid response data"),
			wantErr:  true,
		},
		{
			name:     "Successfully unsmarshal a valid response data format and returns an account data",
			response: loadTestFile("./testdata/fetch_response.json"),
			accountData: AccountData{
				Attributes: &AccountAttributes{
					BankID:       "400300",
					BankIDCode:   "GBDSC",
					BaseCurrency: "GBP",
					Bic:          "NWBKGB22",
					Country:      &[]string{"GB"}[0],
					Name:         []string{"john doe"},
				},
				ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
				OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
				Type:           "accounts",
				Version:        &[]int64{12}[0],
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountData, err := extractAccountDataFromResponse(tt.response)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if accountData != nil {
				assert.IsType(t, &AccountData{}, accountData)
				assert.Equal(t, tt.accountData.ID, accountData.ID)
				assert.Equal(t, tt.accountData.OrganisationID, accountData.OrganisationID)
				assert.Equal(t, tt.accountData.Type, accountData.Type)
				assert.Equal(t, tt.accountData.Version, accountData.Version)

				assert.IsType(t, &AccountAttributes{}, accountData.Attributes)
				assert.Equal(t, tt.accountData.Attributes.BankID, accountData.Attributes.BankID)
				assert.Equal(t, tt.accountData.Attributes.BankIDCode, accountData.Attributes.BankIDCode)
				assert.Equal(t, tt.accountData.Attributes.BaseCurrency, accountData.Attributes.BaseCurrency)
				assert.Equal(t, tt.accountData.Attributes.Bic, accountData.Attributes.Bic)
				assert.Equal(t, tt.accountData.Attributes.Country, accountData.Attributes.Country)
				assert.Equal(t, tt.accountData.Attributes.Name, accountData.Attributes.Name)
			}
		})
	}
}

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

			mock.AssertExpectationsForObjects(t, httpUtilsMock)
		})
	}
}

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
