package integration_tests

import (
	"log"
	"net"
	"net/url"
	"os"
	"testing"
	"time"

	"renatoaraujo/form3-account-api-client/accounts"
	"renatoaraujo/form3-account-api-client/httputils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func TestMain(m *testing.M) {
	parsedUri, err := url.ParseRequestURI(getEnv("API_BASE_URI", "https://api.form3.tech"))
	if err != nil {
		panic("failed to parse the base uri, please check your environment variables")
	}

	log.Println("checking if the host is available")
	timeout := time.Duration(1) * time.Second
	conn, err := net.DialTimeout("tcp", parsedUri.Host, timeout)
	if err != nil {
		log.Println(err)
		log.Println("host unreachable, skipping functional tests")
		os.Exit(0)
	}
	defer conn.Close()

	exitVal := m.Run()
	os.Exit(exitVal)
}

func clientSetup() accounts.Client {
	httpClient, _ := httputils.NewClient(getEnv("API_BASE_URI", "https://api.form3.tech"), 15)
	return accounts.NewClient(httpClient)
}

func createAccountResource(accountID uuid.UUID) (*accounts.AccountData, error) {
	client := clientSetup()
	accountData := getAccountData(accountID)

	return client.CreateResource(accountData)
}

func getAccountData(accountID uuid.UUID) *accounts.AccountData {
	return &accounts.AccountData{
		Attributes: &accounts.AccountAttributes{
			BankID:       "400300",
			BankIDCode:   "GBDSC",
			BaseCurrency: "GBP",
			Bic:          "NWBKGB22",
			Country:      &[]string{"GB"}[0],
			Name:         []string{"john doe"},
		},
		ID:             accountID.String(),
		OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
		Type:           "accounts",
		Version:        0,
	}
}

func TestCreateAccount(t *testing.T) {
	tests := []struct {
		name string
		f    func(t *testing.T)
	}{
		{
			name: "Successfully creates an account",
			f: func(t *testing.T) {
				accountID, err := uuid.NewUUID()
				require.NoError(t, err)

				accountData, err := createAccountResource(accountID)
				require.NoError(t, err)

				assert.Equal(t, accountID.String(), accountData.ID)
			},
		},
		{
			name: "Failed to create duplicated account",
			f: func(t *testing.T) {
				accountID, err := uuid.NewUUID()
				require.NoError(t, err)

				_, err = createAccountResource(accountID)
				require.NoError(t, err)

				_, err = createAccountResource(accountID)
				require.Error(t, err)
			},
		},
		{
			name: "Failed to create with invalid account data",
			f: func(t *testing.T) {
				client := clientSetup()
				accountData := &accounts.AccountData{
					ID: "invalid account id",
				}

				_, err := client.CreateResource(accountData)
				require.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.f)
	}
}

func TestFetchAccount(t *testing.T) {
	tests := []struct {
		name string
		f    func(t *testing.T)
	}{
		{
			name: "Successfully fetches an account",
			f: func(t *testing.T) {
				client := clientSetup()
				accountID, err := uuid.NewUUID()
				require.NoError(t, err)

				createdAccountData, err := createAccountResource(accountID)
				require.NoError(t, err)

				assert.Equal(t, getAccountData(accountID), createdAccountData)

				fetchedAccountData, err := client.FetchResource(accountID)
				assert.Equal(t, getAccountData(accountID), fetchedAccountData)
			},
		},
		{
			name: "Failed to fetch an account with an non existent id",
			f: func(t *testing.T) {
				client := clientSetup()
				accountID, err := uuid.NewUUID()
				require.NoError(t, err)

				_, err = client.FetchResource(accountID)
				require.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.f)
	}
}

func TestDeleteAccount(t *testing.T) {
	tests := []struct {
		name string
		f    func(t *testing.T)
	}{
		{
			name: "Successfully deletes an account",
			f: func(t *testing.T) {
				client := clientSetup()
				accountID, err := uuid.NewUUID()
				require.NoError(t, err)

				createdAccountData, err := createAccountResource(accountID)
				require.NoError(t, err)

				err = client.DeleteResource(accountID, createdAccountData.Version)
				require.NoError(t, err)
			},
		},
		{
			name: "Failed to delete an non existent account",
			f: func(t *testing.T) {
				client := clientSetup()
				accountID, err := uuid.NewUUID()
				require.NoError(t, err)

				err = client.DeleteResource(accountID, 0)
				require.Error(t, err)
				require.EqualError(t, err, "api failure with status code 404 and message: not found; unable to delete resource")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.f)
	}
}
