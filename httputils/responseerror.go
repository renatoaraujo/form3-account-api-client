package httputils

import "fmt"

type ResponseError struct {
	ErrorMessage string `json:"error_message,omitempty"`
	StatusCode   int
}

func (err *ResponseError) Error() string {
	if err.ErrorMessage == "" {
		return fmt.Sprintf("api failure with status code %d and no message received", err.StatusCode)
	}
	return fmt.Sprintf("api failure with status code %d and message: %s", err.StatusCode, err.ErrorMessage)
}
