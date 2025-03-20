package client

import (
	"context"
	"github.com/JacobPotter/go-zendesk/credentialtypes"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

const (
	baseURLFormat      = "https://%s.zendesk.com/api/v2"
	baseSuncoURLFormat = "https://%s.zendesk.com/sc/v2/apps/%s"
)

var defaultHeaders = map[string]string{
	"User-Agent":   "JacobPotter/go-zendesk/0.18.0",
	"Content-Type": "application/json",
}

var subdomainRegexp = regexp.MustCompile("^[a-z0-9][a-z0-9-]+[a-z0-9]$")

type RetryMeta struct {
	ClientRetry bool
	WaitTime    time.Duration
}

type (
	// BaseClient of Zendesk API
	BaseClient struct {
		BaseURL    *url.URL
		HttpClient *http.Client
		Credential credentialtypes.Credential
		Headers    map[string]string
		sunco      bool
		suncoAppId string
		Retry      RetryMeta
	}

	// BaseAPI encapsulates base methods for zendesk client
	BaseAPI interface {
		Get(ctx context.Context, path string) ([]byte, error)
		Post(ctx context.Context, path string, data interface{}) ([]byte, error)
		Put(ctx context.Context, path string, data interface{}) ([]byte, error)
		Delete(ctx context.Context, path string) error
	}

	// CursorPagination contains options for using cursor pagination.
	// Cursor pagination is preferred where possible.
	CursorPagination struct {
		// PageSize sets the number of results per page.
		// Most endpoints support up to 100 records per page.
		PageSize int `url:"page[size],omitempty"`

		// PageAfter provides the "next" cursor.
		PageAfter string `url:"page[after],omitempty"`

		// PageBefore provides the "previous" cursor.
		PageBefore string `url:"page[before],omitempty"`
	}

	// CursorPaginationMeta contains information concerning how to fetch
	// next and previous results, and if next results exist.
	CursorPaginationMeta struct {
		// HasMore is true if more results exist in the endpoint.
		HasMore bool `json:"has_more,omitempty"`

		// AfterCursor contains the cursor of the next result set.
		AfterCursor string `json:"after_cursor,omitempty"`

		// BeforeCursor contains the cursor of the previous result set.
		BeforeCursor string `json:"before_cursor,omitempty"`
	}
)

// NewBaseClient creates new Zendesk API client
func NewBaseClient(httpClient *http.Client, sunco bool) (*BaseClient, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	client := &BaseClient{HttpClient: httpClient, sunco: sunco}
	client.Headers = defaultHeaders
	return client, nil
}
