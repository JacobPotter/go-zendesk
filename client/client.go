package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JacobPotter/go-zendesk/credentialtypes"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
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

type (
	// BaseClient of Zendesk API
	BaseClient struct {
		BaseURL    *url.URL
		HttpClient *http.Client
		Credential credentialtypes.Credential
		Headers    map[string]string
		sunco      bool
		suncoAppId string
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

func (c *BaseClient) SetSuncoAppId(suncoAppId string) {
	c.suncoAppId = suncoAppId
}

// NewBaseClient creates new Zendesk API client
func NewBaseClient(httpClient *http.Client, sunco bool) (*BaseClient, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	client := &BaseClient{HttpClient: httpClient, sunco: sunco}
	client.Headers = defaultHeaders
	return client, nil
}

// SetHeader saves HTTP header in client. It will be included all API request
func (c *BaseClient) SetHeader(key string, value string) {
	c.Headers[key] = value
}

// SetSubdomain saves subdomain in client. It will be used
// when call API
func (c *BaseClient) SetSubdomain(subdomain string) error {
	if !subdomainRegexp.MatchString(subdomain) {
		return fmt.Errorf("%s is invalid subdomain", subdomain)
	}

	var baseURLString string

	if c.sunco {
		if c.suncoAppId != "" {
			log.Fatal("cannot set sunco to true without also setting suncoAppId")
		}
		baseURLString = fmt.Sprintf(baseSuncoURLFormat, subdomain, c.suncoAppId)
	} else {
		baseURLString = fmt.Sprintf(baseURLFormat, subdomain)
	}
	baseURL, err := url.Parse(baseURLString)
	if err != nil {
		return err
	}

	c.BaseURL = baseURL
	return nil
}

// SetEndpointURL replace full URL of endpoint without subdomain validation.
// This is mainly used for testing to point to mock API server.
func (c *BaseClient) SetEndpointURL(newURL string) error {
	baseURL, err := url.Parse(newURL)
	if err != nil {
		return err
	}

	c.BaseURL = baseURL
	return nil
}

// SetCredential saves credential in client. It will be set
// to request header when call API
func (c *BaseClient) SetCredential(cred credentialtypes.Credential) {
	if c.sunco {
		if _, ok := cred.(credentialtypes.BasicAuthCredential); !ok {
			c.Credential = cred
		} else {
			log.Fatal("Invalid Credential Type, Only Basic Auth Credentials allowed")
		}
	} else {
		c.Credential = cred
	}
}

// Get Get JSON data from API and returns its body as []bytes
func (c *BaseClient) Get(ctx context.Context, path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.BaseURL.String()+path, nil)
	if err != nil {
		return nil, err
	}

	req = c.PrepareRequest(ctx, req)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		waitForRetry(resp)
		return c.Get(ctx, path)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, Error{
			ErrorBody: body,
			Resp:      resp,
		}
	}
	return body, nil
}

func waitForRetry(resp *http.Response) {
	var retry int64
	retry, err := strconv.ParseInt(resp.Header.Get("retry-after"), 10, 64)
	if err != nil {
		fmt.Printf("Error getting retry header, trying secondary header")
		retry, err = strconv.ParseInt(resp.Header.Get("ratelimit-reset"), 10, 64)
		if err != nil {
			fmt.Printf("Error getting retry header, setting retry after 60 sec by default")
			retry = 60
		}
	}
	time.Sleep(time.Duration(int64(time.Second)*retry + 1))
}

// Post send data to API and returns response body as []bytes
func (c *BaseClient) Post(ctx context.Context, path string, data interface{}) ([]byte, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.BaseURL.String()+path, strings.NewReader(string(bytes)))
	if err != nil {
		return nil, err
	}

	req = c.PrepareRequest(ctx, req)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		waitForRetry(resp)
		return c.Post(ctx, path, data)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
		return nil, Error{
			ErrorBody: body,
			Resp:      resp,
		}
	}

	return body, nil
}

// Put sends data to API and returns response body as []bytes
func (c *BaseClient) Put(ctx context.Context, path string, data interface{}) ([]byte, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, c.BaseURL.String()+path, strings.NewReader(string(bytes)))
	if err != nil {
		return nil, err
	}

	req = c.PrepareRequest(ctx, req)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		waitForRetry(resp)
		return c.Put(ctx, path, data)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// NOTE: some webhook mutation APIs return status No Content.
	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent) {
		return nil, Error{
			ErrorBody: body,
			Resp:      resp,
		}
	}

	return body, nil
}

// Patch sends data to API and returns response body as []bytes
func (c *BaseClient) Patch(ctx context.Context, path string, data interface{}) ([]byte, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPatch, c.BaseURL.String()+path, strings.NewReader(string(bytes)))
	if err != nil {
		return nil, err
	}

	req = c.PrepareRequest(ctx, req)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		waitForRetry(resp)
		return c.Patch(ctx, path, data)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// NOTE: some webhook mutation APIs return status No Content.
	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent) {
		return nil, Error{
			ErrorBody: body,
			Resp:      resp,
		}
	}

	return body, nil
}

// Delete sends data to API and returns an error if unsuccessful
func (c *BaseClient) Delete(ctx context.Context, path string) error {
	req, err := http.NewRequest(http.MethodDelete, c.BaseURL.String()+path, nil)
	if err != nil {
		return err
	}

	req = c.PrepareRequest(ctx, req)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		waitForRetry(resp)
		return c.Delete(ctx, path)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return Error{
			ErrorBody: body,
			Resp:      resp,
		}
	}

	return nil
}

// prepare request sets common request variables such as authn and user agent
func (c *BaseClient) PrepareRequest(ctx context.Context, req *http.Request) *http.Request {
	out := req.WithContext(ctx)
	c.IncludeHeaders(out)
	if c.Credential != nil {
		if c.Credential.Bearer() {
			out.Header.Add("Authorization", "Bearer "+c.Credential.Secret())
		} else {
			out.SetBasicAuth(c.Credential.Email(), c.Credential.Secret())
		}
	}

	return out
}

// IncludeHeaders set HTTP Headers from client.headers to *http.Request
func (c *BaseClient) IncludeHeaders(req *http.Request) {
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}
}

// AddOptions build query string
func AddOptions(s string, opts interface{}) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// GetData is a generic helper function that retrieves and unmarshals JSON data from a specified URL.
// It takes four parameters:
// - a pointer to a BaseClient (z) which is used to execute the GET request,
// - a context (ctx) for managing the request's lifecycle,
// - a string (url) representing the endpoint from which data should be retrieved,
// - and an empty interface (data) where the retrieved data will be stored after being unmarshalled from JSON.
//
// The function starts by sending a GET request to the specified URL. If the request is successful,
// the returned body in the form of a byte slice is unmarshalled into the provided empty interface using the json.Unmarshal function.
//
// If an error occurs during either the GET request or the JSON unmarshalling, the function will return this error.
func GetData(z BaseAPI, ctx context.Context, url string, data any) error {
	body, err := z.Get(ctx, url)
	if err == nil {
		err = json.Unmarshal(body, data)
		if err != nil {
			return err
		}
	}
	return err
}
