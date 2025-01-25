package zendesk

import (
	"errors"
	"fmt"
	"github.com/JacobPotter/go-zendesk/internal/client"
	"github.com/JacobPotter/go-zendesk/internal/testhelper"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

////////// Helper //////////

func NewTestBaseClient(mockAPI *httptest.Server) *client.BaseClient {
	c := &client.BaseClient{
		HttpClient: http.DefaultClient,
		Credential: client.NewAPITokenCredential("", ""),
	}
	err := c.SetEndpointURL(mockAPI.URL)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
	return c
}

////////// Test //////////

func TestNewClient(t *testing.T) {
	if _, err := client.NewBaseClient(nil, false); err != nil {
		t.Fatal("Failed to create BaseClient")
	}
}

func TestSetHeader(t *testing.T) {
	c, _ := client.NewBaseClient(nil, false)
	c.SetHeader("Header1", "hogehoge")

	if c.Headers["Header1"] != "hogehoge" {
		t.Fatal("Header1 is wrong")
	}
}

func TestSetSubdomainSuccess(t *testing.T) {
	validSubdomain := "subdomain"

	c, _ := client.NewBaseClient(&http.Client{}, false)
	if err := c.SetSubdomain(validSubdomain); err != nil {
		t.Fatal("SetSubdomain should success")
	}
}

func TestSetSubdomainFail(t *testing.T) {
	invalidSubdomain := ".subdomain"

	c, _ := client.NewBaseClient(&http.Client{}, false)
	if err := c.SetSubdomain(invalidSubdomain); err == nil {
		t.Fatal("SetSubdomain should fail")
	}
}

func TestSetEndpointURL(t *testing.T) {
	c, _ := client.NewBaseClient(nil, false)
	if err := c.SetEndpointURL("http://127.0.0.1:3000"); err != nil {
		t.Fatal("SetEndpointURL should success")
	}
}

func TestSetCredential(t *testing.T) {
	c, _ := client.NewBaseClient(nil, false)
	cred := client.NewBasicAuthCredential("john.doe@example.com", "password")
	c.SetCredential(cred)

	if email := c.Credential.Email(); email != "john.doe@example.com" {
		t.Fatal("client.credential.Email() returns wrong email: " + email)
	}
	if secret := c.Credential.Secret(); secret != "password" {
		t.Fatal("client.credential.Secret() returns wrong secret: " + secret)
	}
}

func TestBearerAuthCredential(t *testing.T) {
	c, _ := client.NewBaseClient(nil, false)
	cred := client.NewBearerTokenCredential("hello")
	c.SetCredential(cred)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer hello" {
			t.Fatal("unexpected auth header: " + auth)
		}
	}))
	err := c.SetEndpointURL(server.URL)
	if err != nil {
		t.Logf("Error: %s", err.Error())
	}
	defer server.Close()

	// trigger request, assert in the server code
	_, _ = c.Get(ctx, "/groups.json")
}

func TestGet(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "groups.json")
	c := NewTestBaseClient(mockAPI)
	defer mockAPI.Close()

	body, err := c.Get(ctx, "/groups.json")
	if err != nil {
		t.Fatalf("Failed to send request: %s", err)
	}

	if len(body) == 0 {
		t.Fatal("Response body is empty")
	}
}

func TestGetData(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "groups.json")
	c := NewTestBaseClient(mockAPI)
	defer mockAPI.Close()

	var data struct {
		Groups []Group `json:"groups"`
		Page
	}

	opts := &OBPOptions{}

	u, err := client.AddOptions("/groups.json", opts)
	if err != nil {
		t.Fatal(err)
	}

	err = client.GetData(c, ctx, u, &data)
	if err != nil {
		t.Fatalf("Failed to send request: %s", err)
	}
	if len(data.Groups) == 0 {
		t.Fatal("Response body is empty")
	}
}

func TestGetFailure(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodGet, "groups.json", http.StatusInternalServerError)
	c := NewTestBaseClient(mockAPI)
	defer mockAPI.Close()

	_, err := c.Get(ctx, "/groups.json")
	if err == nil {
		t.Fatal("Did not receive error from client")
	}

	var e client.Error
	if !errors.As(err, &e) {
		t.Fatalf("Did not return a zendesk error %s", err)
	}
}

func TestGetFailureCanReadErrorBody(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodGet, "groups.json", http.StatusInternalServerError)
	c := NewTestBaseClient(mockAPI)
	defer mockAPI.Close()

	_, err := c.Get(ctx, "/groups.json")
	if err == nil {
		t.Fatal("Did not receive error from client")
	}

	var clientErr client.Error
	if ok := errors.As(err, &clientErr); !ok {
		t.Fatalf("Did not return a zendesk error %s", err)
	}

	body := clientErr.Body()
	_, err = io.ReadAll(body)
	if err != nil {
		t.Fatal("BaseClient received error while reading client body")
	}

	err = body.Close()
	if err != nil {
		t.Fatal("BaseClient received error while closing body")
	}
}

func TestPost(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPost, "groups.json", http.StatusCreated)
	c := NewTestBaseClient(mockAPI)
	defer mockAPI.Close()

	body, err := c.Post(ctx, "/groups.json", Group{})
	if err != nil {
		t.Fatalf("Failed to send request: %s", err)
	}

	if len(body) == 0 {
		t.Fatal("Response body is empty")
	}
}

func TestPostFailure(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPost, "groups.json", http.StatusInternalServerError)
	c := NewTestBaseClient(mockAPI)
	defer mockAPI.Close()

	_, err := c.Post(ctx, "/groups.json", Group{})
	if err == nil {
		t.Fatal("Did not receive error from client")
	}

	var e client.Error
	if !errors.As(err, &e) {
		t.Fatalf("Did not return a zendesk error %s", err)
	}
}

func TestPut(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPut, "groups.json", http.StatusOK)
	c := NewTestBaseClient(mockAPI)
	defer mockAPI.Close()

	body, err := c.Put(ctx, "/groups.json", Group{})
	if err != nil {
		t.Fatalf("Failed to send request: %s", err)
	}

	if len(body) == 0 {
		t.Fatal("Response body is empty")
	}
}

func TestPutFailure(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPut, "groups.json", http.StatusInternalServerError)
	c := NewTestBaseClient(mockAPI)
	defer mockAPI.Close()

	_, err := c.Put(ctx, "/groups.json", Group{})
	if err == nil {
		t.Fatal("Did not receive error from client")
	}

	var clientErr client.Error
	if !errors.As(err, &clientErr) {
		t.Fatalf("Did not return a zendesk error %s", err)
	}
}

func TestDelete(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestBaseClient(mockAPI)
	err := c.Delete(ctx, "/foo/id")
	if err != nil {
		t.Fatalf("Failed to send request: %s", err)
	}
}

func TestDeleteFailure(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestBaseClient(mockAPI)
	err := c.Delete(ctx, "/foo/id")
	if err == nil {
		t.Fatalf("Failed to recieve error from Delete")
	}
}

func TestIncludeHeaders(t *testing.T) {
	c, _ := client.NewBaseClient(nil, false)
	c.Headers = map[string]string{
		"Header1":      "1",
		"Header2":      "2",
		"Content-Type": "application/json",
	}

	req, _ := http.NewRequest("POST", "localhost", strings.NewReader(""))
	c.IncludeHeaders(req)

	if len(req.Header) != 3 {
		t.Fatal("req.Header length does not match")
	}

	for k, v := range req.Header {
		switch k {
		case "Header1":
			if v[0] != "1" {
				t.Fatalf(`%s header expect "1", but got "%s"`, k, v[0])
			}
		case "Header2":
			if v[0] != "2" {
				t.Fatalf(`%s header expect "2", but got "%s"`, k, v[0])
			}
		case "Content-Type":
			if v[0] != "application/json" {
				t.Fatalf(`%s header expect "2", but got "%s"`, k, v[0])
			}
		}
	}
}

func TestAddOptions(t *testing.T) {
	ep := "/triggers.json"
	ops := &TriggerListOptions{
		PageOptions: PageOptions{
			PerPage: 10,
			Page:    2,
		},
		Active: true,
	}
	expected := "/triggers.json?active=true&page=2&per_page=10"

	u, err := client.AddOptions(ep, ops)
	if err != nil {
		t.Fatal(err)
	}

	if u != expected {
		t.Fatalf("\nExpect:\t%s\nGot:\t%s", expected, u)
	}
}
