package zendesk

import (
	"fmt"
	"github.com/JacobPotter/go-zendesk/testhelper"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestGetUserFields(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "user_fields.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	fields, page, err := c.GetUserFields(ctx, nil)
	if err != nil {
		t.Fatalf("Received error calling API: %v", err)
	}

	if page.Count != 1 {
		t.Fatalf("Did not receive the correct count in the page field. Was %d expected 1", page.Count)
	}

	n := len(fields)
	if n != 1 {
		t.Fatalf("Expected 1 entry in fields list. Got %d", n)
	}

	id := fields[0].ID
	if id != 7 {
		t.Fatalf("Field did not have the expected id. Was %d", id)
	}
}

func TestUserFieldQueryParamsSet(t *testing.T) {
	opts := UserFieldListOptions{
		PageOptions{
			Page: 2,
		},
	}
	expected := fmt.Sprintf("page=%d", opts.Page)
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryString := r.URL.Query().Encode()
		if queryString != expected {
			t.Fatalf(`Did not get the expect query string: "%s". Was: "%s"`, expected, queryString)
		}
		_, err := w.Write(testhelper.ReadFixture(t, filepath.Join(http.MethodGet, "user_fields.json")))
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	defer mockAPI.Close()
	c := NewTestClient(mockAPI)
	_, _, err := c.GetUserFields(ctx, &opts)
	if err != nil {
		t.Fatalf("Received error calling API: %v", err)
	}
}
