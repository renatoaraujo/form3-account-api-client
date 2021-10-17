package accounts

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
)

const basePath = "/v1/organisation/accounts"

type httpUtils interface {
	Delete(resourcePath string, query map[string]string) error
	Get(resourcePath string) ([]byte, error)
	Post(resourcePath string, body []byte) ([]byte, error)
}

// Client is the representation of the client to interact with the account section on form3 api see https://api-docs.form3.tech/api.html#organisation-accounts
type Client struct {
	http httpUtils
}

// NewClient creates a new account client instance with a http utils
func NewClient(httpUtils httpUtils) Client {
	return Client{http: httpUtils}
}

func extractAccountDataFromResponse(response []byte) (*AccountData, error) {
	responsePayload := &payload{}
	if err := json.Unmarshal(response, responsePayload); err != nil {
		return nil, errors.New("failed to unmarshal response data")
	}

	return responsePayload.Data, nil
}

// CreateResource creates a new account resource see https://api-docs.form3.tech/api.html#organisation-accounts-create
func (client *Client) CreateResource(accountData *AccountData) (*AccountData, error) {
	requestPayload, err := json.Marshal(&payload{
		Data: accountData,
	})
	if err != nil {
		return nil, fmt.Errorf("%w; unable to convert account data payload", err)
	}

	response, err := client.http.Post(basePath, requestPayload)
	if err != nil {
		return nil, fmt.Errorf("%w; unable to create resource", err)
	}

	responseAccountData, err := extractAccountDataFromResponse(response)
	if err != nil {
		return nil, fmt.Errorf("%w; failed to extract account data after successful account creation", err)
	}

	return responseAccountData, nil
}

// FetchResource fetches an account resource by an account id see https://api-docs.form3.tech/api.html#organisation-accounts-fetch
func (client *Client) FetchResource(accountID uuid.UUID) (*AccountData, error) {
	resourcePath := fmt.Sprintf("%s/%s", basePath, accountID.String())
	response, err := client.http.Get(resourcePath)
	if err != nil {
		return nil, fmt.Errorf("%w; unable to fetch resource", err)
	}

	responseAccountData, err := extractAccountDataFromResponse(response)
	if err != nil {
		return nil, fmt.Errorf("%w; unable to extract the fetched data from the response", err)
	}

	return responseAccountData, nil
}

// DeleteResource deletes an account resource by an account id and version see https://api-docs.form3.tech/api.html#organisation-accounts-delete
func (client *Client) DeleteResource(accountID uuid.UUID, version int) error {
	resourcePath := fmt.Sprintf("%s/%s", basePath, accountID.String())
	query := map[string]string{
		"version": strconv.Itoa(version),
	}
	if err := client.http.Delete(resourcePath, query); err != nil {
		return fmt.Errorf("%w; unable to delete resource", err)
	}

	return nil
}
