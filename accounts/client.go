package accounts

import (
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
