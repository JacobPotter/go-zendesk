package zendesk

import (
	"fmt"
	"github.com/JacobPotter/go-zendesk/internal/client"
	"net/http"
	"net/http/httptest"
)

func NewTestClient(mockAPI *httptest.Server) *Client {
	c := &Client{
		BaseClient: &client.BaseClient{HttpClient: http.DefaultClient,
			Credential: client.NewAPITokenCredential("", "")},
	}
	err := c.SetEndpointURL(mockAPI.URL)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
	return c
}
