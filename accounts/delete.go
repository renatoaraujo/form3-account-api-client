package accounts

import (
	"fmt"
)

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
