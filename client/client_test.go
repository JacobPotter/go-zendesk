package client

import (
	"fmt"
	"github.com/JacobPotter/go-zendesk/credentialtypes"
	"net/http"
	"net/http/httptest"
)

func NewTestClient(mockAPI *httptest.Server, clientRetry bool) *BaseClient {

	c := &BaseClient{
		HttpClient:  http.DefaultClient,
		Credential:  credentialtypes.NewAPITokenCredential("", ""),
		ClientRetry: clientRetry,
	}

	err := c.SetEndpointURL(mockAPI.URL)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
	return c
}
