package httputils

type responseError struct {
	ErrorMessage string `json:"error_message,omitempty"`
}
