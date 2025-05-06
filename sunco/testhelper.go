package sunco

import (
	"log"
	"net/http/httptest"
)

func NewTestClient(mockAPI *httptest.Server) *Client {
	c, err := NewClient(nil)
	if err != nil {
		log.Fatalf("Error: %s\n", err.Error())
	}
	err = c.SetEndpointURL(mockAPI.URL)
	if err != nil {
		log.Fatalf("Error: %s\n", err.Error())
	}
	return c
}
