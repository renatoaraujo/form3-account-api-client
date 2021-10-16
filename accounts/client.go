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

func (client *Client) CreateResource(accountData *AccountData) (*AccountData, error) {
	requestPayload, err := json.Marshal(&payload{
		Data: accountData,
	})
	if err != nil {
		return nil, fmt.Errorf("%w; unable to convert account data payload", err)
	}

	apiResponse, err := client.http.Post(basePath, requestPayload)
	if err != nil {
		return nil, fmt.Errorf("%w; unable to create resource", err)
	}

	responsePayload := &payload{}
	if err = json.Unmarshal(apiResponse, responsePayload); err != nil {
		return nil, errors.New("the response from the api was successfully but failed to unmarshal response data")
	}

	return responsePayload.Data, nil
}
