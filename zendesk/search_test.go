package zendesk

import (
	"github.com/JacobPotter/go-zendesk/internal/testhelper"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestSearchTickets(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "search_ticket.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	results, _, err := c.Search(ctx, &SearchOptions{})
	if err != nil {
		t.Fatalf("Failed to get search results: %s", err)
	}

	list := results.List()
	if len(list) != 1 {
		t.Fatalf("expected length of sla policies is , but got %d", len(list))
	}

	ticket, ok := list[0].(Ticket)
	if !ok {
		t.Fatalf("Cannot assert %v as a ticket", list[0])
	}

	if ticket.ID != 4 {
		t.Fatalf("Ticket did not have the expected id %v", ticket)
	}
}

func TestCountTickets(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "search_count_ticket.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	count, err := c.SearchCount(ctx, &CountOptions{})
	if err != nil {
		t.Fatalf("Failed to get count: %s", err)
	}

	expected := 10
	if count != expected {
		t.Fatalf("expected count of tickets is %d, but got %d", expected, count)
	}
}

func TestSearchGroup(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "search_group.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	results, _, err := c.Search(ctx, &SearchOptions{})
	if err != nil {
		t.Fatalf("Failed to get search results: %s", err)
	}

	list := results.List()
	if len(list) != 1 {
		t.Fatalf("expected length of sla policies is , but got %d", len(list))
	}

	result, ok := list[0].(Group)
	if !ok {
		t.Fatalf("Cannot assert %v as a group", list[0])
	}

	if result.ID != 360007194452 {
		t.Fatalf("Group did not have the expected id %v", result)
	}
}

func TestSearchUser(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "search_user.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	results, _, err := c.Search(ctx, &SearchOptions{})
	if err != nil {
		t.Fatalf("Failed to get search results: %s", err)
	}

	list := results.List()
	if len(list) != 1 {
		t.Fatalf("expected length of sla policies is , but got %d", len(list))
	}

	result, ok := list[0].(User)
	if !ok {
		t.Fatalf("Cannot assert %v as a user", list[0])
	}

	if result.ID != 1234 {
		t.Fatalf("Group did not have the expected id %v", result)
	}
}

func TestSearchQueryParam(t *testing.T) {
	expected := "query string"
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryString := r.URL.Query().Get("query")
		if queryString != expected {
			t.Fatalf(`Did not get the expect query string: "%s". Was: "%s"`, expected, queryString)
		}
		_, err := w.Write(testhelper.ReadFixture(t, filepath.Join(http.MethodGet, "search_user.json")))
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	opts := SearchOptions{
		PageOptions: PageOptions{
			Page:    1,
			PerPage: 2,
		},
		Query: expected,
	}

	_, _, err := c.Search(ctx, &opts)
	if err != nil {
		t.Fatalf("Received error from search api")
	}
}
