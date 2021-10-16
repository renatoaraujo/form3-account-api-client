package accounts

import (
	"encoding/json"
	"errors"
	"fmt"
)

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
