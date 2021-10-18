package integration_tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	log.Println("checking if the host is available, this is to prevent running the tests without running the docker")
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

func createAccountResource(accountData *accounts.AccountData) (*accounts.AccountData, error) {
	client := clientSetup()

	return client.CreateResource(accountData)
}

func getCreateAccountData(accountID uuid.UUID) *accounts.AccountData {
	accountData := loadAccountDataFromFileWithCustomID("./testdata/account_create_data.json", accountID)

	return accountData
}

func getFetchAccountData(accountID uuid.UUID) *accounts.AccountData {
	accountData := loadAccountDataFromFileWithCustomID("./testdata/account_fetch_data.json", accountID)

	return accountData
}

func loadAccountDataFromFileWithCustomID(file string, accountID uuid.UUID) *accounts.AccountData {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		panic("failed to load the test data file")
	}

	var payload accounts.Payload
	if err = json.Unmarshal(raw, &payload); err != nil {
		panic("failed to unmarshal the test data file")
	}

	payload.Data.ID = accountID.String()

	return payload.Data
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

				expectedAccountData := getCreateAccountData(accountID)
				accountData, err := createAccountResource(expectedAccountData)
				require.NoError(t, err)

				assert.Equal(t, expectedAccountData, accountData)
			},
		},
		{
			name: "Failed to create duplicated account",
			f: func(t *testing.T) {
				accountID, err := uuid.NewUUID()
				require.NoError(t, err)

				_, err = createAccountResource(getCreateAccountData(accountID))
				require.NoError(t, err)

				_, err = createAccountResource(getCreateAccountData(accountID))
				require.Error(t, err)
				require.EqualError(t, err,
					"api failure with status code 409 and message: Account cannot be created as it violates a duplicate constraint; unable to create resource",
				)
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

				_, err = createAccountResource(getCreateAccountData(accountID))
				require.NoError(t, err)

				actual, err := client.FetchResource(accountID)
				expected := getFetchAccountData(accountID)

				assert.Equal(t, expected.ID, actual.ID)
				assert.Equal(t, expected.OrganisationID, actual.OrganisationID)
				assert.Equal(t, expected.Type, actual.Type)
				assert.Equal(t, expected.Version, actual.Version)
				assert.Equal(t, expected.Attributes.AccountClassification, actual.Attributes.AccountClassification)
				assert.Equal(t, expected.Attributes.AccountMatchingOptOut, actual.Attributes.AccountMatchingOptOut)
				assert.Equal(t, expected.Attributes.AccountNumber, actual.Attributes.AccountNumber)
				assert.Equal(t, expected.Attributes.AccountQualifier, actual.Attributes.AccountQualifier)
				assert.Equal(t, expected.Attributes.AlternativeNames, actual.Attributes.AlternativeNames)
				assert.Equal(t, expected.Attributes.BankID, actual.Attributes.BankID)
				assert.Equal(t, expected.Attributes.BankIDCode, actual.Attributes.BankIDCode)
				assert.Equal(t, expected.Attributes.BaseCurrency, actual.Attributes.BaseCurrency)
				assert.Equal(t, expected.Attributes.Bic, actual.Attributes.Bic)
				assert.Equal(t, expected.Attributes.CustomerID, actual.Attributes.CustomerID)
				assert.Equal(t, expected.Attributes.Country, actual.Attributes.Country)
				assert.Equal(t, expected.Attributes.Iban, actual.Attributes.Iban)
				assert.Equal(t, expected.Attributes.JointAccount, actual.Attributes.JointAccount)
				assert.Equal(t, expected.Attributes.Name, actual.Attributes.Name)
				assert.Equal(t, expected.Attributes.ProcessingService, actual.Attributes.ProcessingService)
				assert.Equal(t, expected.Attributes.ReferenceMask, actual.Attributes.ReferenceMask)
				assert.Equal(t, expected.Attributes.SecondaryIdentification, actual.Attributes.SecondaryIdentification)
				assert.Equal(t, expected.Attributes.Switched, actual.Attributes.Switched)
				assert.Equal(t, expected.Attributes.UserDefinedInformation, actual.Attributes.UserDefinedInformation)
				assert.Equal(t, expected.Attributes.ValidationType, actual.Attributes.ValidationType)
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
				require.EqualError(t, err,
					fmt.Sprintf("api failure with status code 404 and message: record %s does not exist; unable to fetch resource", accountID.String()),
				)
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

				accountData := getCreateAccountData(accountID)
				createdAccountData, err := createAccountResource(accountData)
				require.NoError(t, err)

				err = client.DeleteResource(accountID, createdAccountData.Version)
				require.NoError(t, err)

				_, err = client.FetchResource(accountID)
				require.Error(t, err)
				require.EqualError(t, err,
					fmt.Sprintf("api failure with status code 404 and message: record %s does not exist; unable to fetch resource", accountID.String()),
				)
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
