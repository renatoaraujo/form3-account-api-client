package accounts

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (client *Client) DeleteResource(accountID string, version int) error {
	_, err := uuid.Parse(accountID)
	if err != nil {
		return errors.New(
			fmt.Sprintf("invalid uuid format: %s", accountID),
		)
	}

	resourcePath := fmt.Sprintf("%s/%s?version=%d", basePath, accountID, version)
	err = client.http.Delete(resourcePath)
	if err != nil {
		return err
	}

	return nil
}
