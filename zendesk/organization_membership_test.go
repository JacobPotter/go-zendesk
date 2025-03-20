package zendesk

import (
	"github.com/JacobPotter/go-zendesk/internal/testhelper"
	"net/http"
	"testing"
)

func TestGetOrganizationMemberships(t *testing.T) {
	mockAPI := testhelper.NewMockAPI(t, http.MethodGet, "organization_memberships.json")
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	orgMemberships, _, err := c.GetOrganizationMemberships(ctx, &OrganizationMembershipListOptions{})
	if err != nil {
		t.Fatalf("Failed to get organization memberships: %s", err)
	}

	expectedOrgMemberships := 2

	if len(orgMemberships) != expectedOrgMemberships {
		t.Fatalf("expected length of organization memberships is %d, but got %d", expectedOrgMemberships, len(orgMemberships))
	}
}

func TestCreateOrganizationMembership(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPost, "organization_membership.json", http.StatusCreated, nil, false)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	_, err := c.CreateOrganizationMembership(ctx, OrganizationMembershipOptions{})
	if err != nil {
		t.Fatalf("Failed to send request to create organization membership: %s", err)
	}
}

func TestSetDefaultOrganization(t *testing.T) {
	mockAPI := testhelper.NewMockAPIWithStatus(t, http.MethodPut, "organization_membership.json", http.StatusOK, nil, false)
	c := NewTestClient(mockAPI)
	defer mockAPI.Close()

	orgMembership, err := c.SetDefaultOrganization(ctx, OrganizationMembershipOptions{})
	if err != nil {
		t.Fatalf("Failed to set the default organization for user: %s", err)
	}

	expectedDefault := true
	if orgMembership.Default != expectedDefault {
		t.Fatalf("Returned org membership does not have the expected default status %v. It is %v", expectedDefault, orgMembership.Default)
	}
}
