package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// Error an error type containing the http response from zendesk
type Error struct {
	ErrorBody []byte
	Resp      *http.Response
}

// NewError is a function to initialize the Error type. This function will be useful
// for unit testing and mocking purposes in the client side
// to test their behavior by the API response.
func NewError(body []byte, resp *http.Response) Error {
	return Error{
		ErrorBody: body,
		Resp:      resp,
	}
}

// Error the error string for this error
func (e Error) Error() string {
	msg := string(e.ErrorBody)
	if msg == "" {
		msg = http.StatusText(e.Status())
	}

	return fmt.Sprintf("%d: %s", e.Resp.StatusCode, msg)
}

// Body is the Body of the HTTP response
func (e Error) Body() io.ReadCloser {
	return io.NopCloser(bytes.NewBuffer(e.ErrorBody))
}

// Headers the HTTP headers returned from zendesk
func (e Error) Headers() http.Header {
	return e.Resp.Header
}

// Status the HTTP status code returned from zendesk
func (e Error) Status() int {
	return e.Resp.StatusCode
}

// OptionsError is an error type for invalid option argument.
type OptionsError struct {
	Opts interface{}
}

func (e *OptionsError) Error() string {
	return fmt.Sprintf("invalid options: %v", e.Opts)
}
