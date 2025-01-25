package sunco

import (
	"github.com/JacobPotter/go-zendesk/internal/client"
	"net/http"
)

// API an interface containing all the zendesk client methods
type API interface {
	client.BaseAPI
}

type Client struct {
	*client.BaseClient
}

func NewClient(httpClient *http.Client) (*Client, error) {
	suncoClient, err := client.NewBaseClient(httpClient, true)
	if err != nil {
		return nil, err
	}
	return &Client{BaseClient: suncoClient}, nil
}

var _ API = (*Client)(nil)
