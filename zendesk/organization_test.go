package zendesk

import (
	"github.com/JacobPotter/go-zendesk/internal/testhelper"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateOrganization(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPost, "organization.json", http.StatusCreated, nil, false)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	_, err := c.CreateOrganization(ctx, Organization{})
	if err != nil {
		t.Fatalf("Failed to send request to create organization: %s", err)
	}
}

func TestGetOrganization(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "organization.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	org, err := c.GetOrganization(ctx, 123)
	if err != nil {
		t.Fatalf("Failed to get organization: %s", err)
	}

	expectedID := int64(361898904439)
	if org.ID != expectedID {
		t.Fatalf("Returned organization does not have the expected ID %d. Organization ID is %d", expectedID, org.ID)
	}
}

func TestGetOrganizations(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "organizations.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	orgs, _, err := c.GetOrganizations(ctx, &OrganizationListOptions{})
	if err != nil {
		t.Fatalf("Failed to get organizations: %s", err)
	}

	if len(orgs) != 2 {
		t.Fatalf("expected length of organizations is , but got %d", len(orgs))
	}
}

func TestUpdateOrganization(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPut, "organization.json", http.StatusOK, nil, false)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	updatedOrg, err := c.UpdateOrganization(ctx, int64(1234), Organization{})
	if err != nil {
		t.Fatalf("Failed to send request to create organization: %s", err)
	}

	expectedID := int64(361898904439)
	if updatedOrg.ID != expectedID {
		t.Fatalf("Updated organization %v did not have expected id %d", updatedOrg, expectedID)
	}
}

func TestDeleteOrganization(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		_, err := w.Write(nil)
		if err != nil {
			t.Logf("Error: %s", err.Error())
		}
	}))

	c := NewTestClient(mockAPI)
	err := c.DeleteOrganization(ctx, 1234)
	if err != nil {
		t.Fatalf("Failed to delete organization: %s", err)
	}
}
