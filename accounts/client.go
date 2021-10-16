package accounts

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const (
	basePath = "/v1/organisation/accounts"
)

type httpUtils interface {
	Delete(resourcePath string) error
	Get(resourcePath string) ([]byte, error)
	Post(resourcePath string, payload []byte) ([]byte, error)
}

type Client struct {
	http httpUtils
}

func NewClient(httpUtils httpUtils) Client {
	return Client{http: httpUtils}
}

func validateAccountIDFormat(accountID string) error {
	_, err := uuid.Parse(accountID)
	if err != nil {
		return errors.New(
			fmt.Sprintf("invalid account id uuid format: %s", accountID),
		)
	}

	return nil
}

func extractAccountDataFromResponse(response []byte) (*AccountData, error) {
	responsePayload := &payload{}
	if err := json.Unmarshal(response, responsePayload); err != nil {
		return nil, errors.New("failed to unmarshal response data")
	}

	return responsePayload.Data, nil
}

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

func (client *Client) FetchResource(accountID string) (*AccountData, error) {
	if err := validateAccountIDFormat(accountID); err != nil {
		return nil, fmt.Errorf("%w; unable to fetch resource", err)
	}

	resourcePath := fmt.Sprintf("%s/%s", basePath, accountID)
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

func (client *Client) DeleteResource(accountID string, version int) error {
	if err := validateAccountIDFormat(accountID); err != nil {
		return fmt.Errorf("%w; unable to delete resource", err)
	}

	resourcePath := fmt.Sprintf("%s/%s?version=%d", basePath, accountID, version)
	if err := client.http.Delete(resourcePath); err != nil {
		return fmt.Errorf("%w; unable to delete resource", err)
	}

	return nil
}
