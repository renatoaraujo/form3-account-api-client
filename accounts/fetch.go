package accounts

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (client *Client) FetchResource(accountID string) (*AccountData, error) {
	_, err := uuid.Parse(accountID)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("invalid uuid: %s", accountID),
		)
	}

	resourcePath := fmt.Sprintf("%s/%s", basePath, accountID)
	apiResponse, err := client.http.Get(resourcePath)
	if err != nil {
		return nil, err
	}

	responseData := &response{}
	err = json.Unmarshal(apiResponse, responseData)
	if err != nil {
		errors.New("failed to unmarshal response data")
	}

	return responseData.Data, nil
}
