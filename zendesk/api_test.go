package zendesk

import (
	"fmt"
	"github.com/JacobPotter/go-zendesk/client"
	"github.com/JacobPotter/go-zendesk/credentialtypes"
	"net/http"
	"net/http/httptest"
)

func NewTestClient(mockAPI *httptest.Server) *Client {
	c := &Client{
		BaseClient: &client.BaseClient{HttpClient: http.DefaultClient,
			Credential: credentialtypes.NewAPITokenCredential("", "")},
	}
	err := c.SetEndpointURL(mockAPI.URL)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
	return c
}
