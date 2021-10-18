package httputils

import "fmt"

// ResponseError is the representation of an error coming from the form3 api with the status code
type ResponseError struct {
	ErrorMessage string `json:"error_message,omitempty"`
	StatusCode   int
}

func (err *ResponseError) Error() string {
	return fmt.Sprintf("api failure with status code %d and message: %s", err.StatusCode, err.ErrorMessage)
}
