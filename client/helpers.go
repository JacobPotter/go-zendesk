package client

import (
	"context"
	"fmt"
	"github.com/JacobPotter/go-zendesk/credentialtypes"
	"github.com/google/go-querystring/query"
	"log"
	"net/http"
	"net/url"
)

func (c *BaseClient) SetSuncoAppId(suncoAppId string) {
	c.suncoAppId = suncoAppId
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
		if c.suncoAppId == "" {
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

// PrepareRequest prepare request sets common request variables such as authn and user agent
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
