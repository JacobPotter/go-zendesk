package zendesk

import (
	"github.com/JacobPotter/go-zendesk/internal/testhelper"
	"net/http"
	"testing"
)

func TestGetOrganizationFields(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "organization_fields.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	ticketFields, _, err := c.GetOrganizationFields(ctx)
	if err != nil {
		t.Fatalf("Failed to get organization fields: %s", err)
	}

	if len(ticketFields) != 2 {
		t.Fatalf("expected length of organization fields is , but got %d", len(ticketFields))
	}
}

func TestOrganizationField(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPost, "organization_fields.json", http.StatusCreated, nil, false)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	_, err := c.CreateOrganizationField(ctx, OrganizationField{})
	if err != nil {
		t.Fatalf("Failed to send request to create organization field: %s", err)
	}
}
