package accounts

import (
	"encoding/json"
	"errors"
	"fmt"
)

func (client *Client) FetchResource(accountID string) (*AccountData, error) {
	err := validateAccountIDFormat(accountID)
	if err != nil {
		return nil, fmt.Errorf("%w; unable to fetch resource", err)
	}

	resourcePath := fmt.Sprintf("%s/%s", basePath, accountID)
	apiResponse, err := client.http.Get(resourcePath)
	if err != nil {
		return nil, fmt.Errorf("%w; unable to fetch resource", err)
	}

	responseData := &response{}
	if err = json.Unmarshal(apiResponse, responseData); err != nil {
		return nil, errors.New("failed to unmarshal response data")
	}

	return responseData.Data, nil
}
